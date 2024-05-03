package acme

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"time"
)

type Acme struct {
	email      string
	idn        []Identifier
	challenges []Challenge
	finalize   string
	options    *Options
	csr        []byte
	privateKey *ecdsa.PrivateKey
}

func NewAcme(email string, opts ...Option) (*Acme, error) {
	options, err := newOptions(opts...)
	if err != nil {
		return nil, err
	}
	// 有kid没有私钥，报错
	if options.client.GetKid() != "" && options.client.GetPrivateKey() == nil {
		return nil, ErrPrivateKeyRequired
	}
	// 没有私钥，生成一个
	if options.client.GetPrivateKey() == nil {
		key, err := GenerateKey()
		if err != nil {
			return nil, err
		}
		options.client.SetPrivateKey(key)
	}
	return &Acme{
		email:   email,
		options: options,
	}, nil
}

type Challenge struct {
	Authorization string
	Challenge     string
	Domain        string
	Record        string
}

func (a *Acme) NewAccount(ctx context.Context) (string, error) {
	account, err := a.options.client.NewAccount(ctx, []string{
		"mailto:" + a.email,
	})
	if err != nil {
		return "", err
	}
	a.options.client.SetKid(account)
	return account, nil
}

func (a *Acme) GetClient() *Client {
	return a.options.client
}

func (a *Acme) GetChallenge(ctx context.Context, domain []string) ([]Challenge, error) {
	if a.options.client.GetKid() == "" {
		// 有kid跳过此步
		account, err := a.options.client.NewAccount(ctx, []string{
			"mailto:" + a.email,
		})
		if err != nil {
			return nil, err
		}
		a.options.client.SetKid(account)
	}
	// 创建订单
	a.idn = make([]Identifier, 0)
	for _, v := range domain {
		a.idn = append(a.idn, Identifier{
			Type:  "dns",
			Value: v,
		})
	}
	auth, err := a.options.client.NewOrder(ctx, a.idn)
	if err != nil {
		return nil, err
	}
	a.challenges = make([]Challenge, 0)
	a.finalize = auth.Finalize
	// 获取授权
	for _, authorization := range auth.Authorizations {
		auth, err := a.options.client.GetAuthorization(ctx, authorization)
		if err != nil {
			return nil, err
		}
		//获取dns-01的验证信息
		if auth.Identifier.Type == "dns" {
			for _, challenge := range auth.Challenges {
				if challenge.Type == "dns-01" {
					a.challenges = append(a.challenges, Challenge{
						Authorization: authorization,
						Challenge:     challenge.Url,
						Domain:        auth.Identifier.Value,
						Record:        a.options.client.DNS01ChallengeRecord(challenge.Token),
					})
				}
			}
			continue
		}
	}
	return a.challenges, nil
}

var ErrChallengeInvalid = errors.New("challenge invalid")

func (a *Acme) WaitChallenge(ctx context.Context) error {
	// 等待所有的验证都完成
	for _, challenge := range a.challenges {
		// 先发送请求验证
		_, err := a.options.client.RequestChallenge(ctx, challenge.Challenge)
		if err != nil {
			return err
		}
		// 轮询验证结果
		for {
			flag := false
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
				auth, err := a.options.client.GetAuthorization(ctx, challenge.Authorization)
				if err != nil {
					return err
				}
				for _, v := range auth.Challenges {
					if v.Type == "dns-01" {
						if v.Status == "valid" {
							flag = true
							break
						}
						if v.Status == "invalid" {
							return ErrChallengeInvalid
						}
					}
				}
				time.Sleep(time.Second * 5)
			}
			if flag {
				break
			}
		}
	}
	return nil
}

func (a *Acme) GetCertificate(ctx context.Context) ([]byte, error) {
	// 创建证书请求
	var err error
	a.csr, a.privateKey, err = CreateCertificateRequest(a.idn)
	if err != nil {
		return nil, err
	}
	// 完成订单
	order, err := a.options.client.Finalize(ctx, a.finalize, a.csr)
	if err != nil {
		return nil, err
	}
	// 获取证书
	cert, err := a.options.client.GetCertificate(ctx, order.Certificate)
	if err != nil {
		return nil, err
	}
	return cert, nil
}

func (a *Acme) GetCSR() []byte {
	return a.csr
}

func (a *Acme) GetPrivateKey() *ecdsa.PrivateKey {
	return a.privateKey
}
