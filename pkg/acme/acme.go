package acme

import (
	"errors"
	"time"
)

type Acme struct {
	email      string
	domain     []string
	idn        []Identifiers
	challenges []Challenge
	finalize   string
	options    *Options
}

func NewAcme(email string, domain []string, opts ...Option) (*Acme, error) {
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
		domain:  domain,
		options: options,
	}, nil
}

type Challenge struct {
	Authorization string
	Challenge     string
	Domain        string
	Record        string
}

func (a *Acme) GetChallenge() error {
	if a.options.client.GetKid() == "" {
		// 有kid跳过此步
		account, err := a.options.client.NewAccount([]string{
			"mailto:" + a.email,
		})
		if err != nil {
			return err
		}
		a.options.client.SetKid(account)
	}
	// 创建订单
	a.idn = make([]Identifiers, 0)
	for _, v := range a.domain {
		a.idn = append(a.idn, Identifiers{
			Type:  "dns",
			Value: v,
		})
	}
	auth, err := a.options.client.NewOrder(a.idn)
	if err != nil {
		return err
	}
	a.challenges = make([]Challenge, 0)
	a.finalize = auth.Finalize
	// 获取授权
	for _, authorization := range auth.Authorizations {
		auth, err := a.options.client.GetAuthorization(authorization)
		if err != nil {
			return err
		}
		//获取dns-01的验证信息
		if auth.Identifier.Type == "dns" {
			for _, challenge := range auth.Challenges {
				if challenge.Type == "dns-01" {
					a.challenges = append(a.challenges, Challenge{
						Authorization: authorization,
						Challenge:     challenge.Url,
						Domain:        auth.Identifier.Value,
						Record:        a.options.client.ChallengeRecord(challenge.Token),
					})
				}
			}
			break
		}
	}
	return nil
}

var ErrChallengeInvalid = errors.New("challenge invalid")

func (a *Acme) WaitChallenge() error {
	// 等待所有的验证都完成
	for _, challenge := range a.challenges {
		// 先发送请求验证
		_, err := a.options.client.RequestChallenge(challenge.Challenge)
		if err != nil {
			return err
		}
		// 轮询验证结果
		for {
			auth, err := a.options.client.GetAuthorization(challenge.Authorization)
			if err != nil {
				return err
			}
			flag := false
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
			if flag {
				break
			}
			time.Sleep(time.Second * 5)
		}
	}
	return nil
}

func (a *Acme) GetCertificate() ([]byte, error) {
	// 创建证书请求
	csr, _, err := CreateCertificateRequest(a.idn)
	if err != nil {
		return nil, err
	}
	// 完成订单
	order, err := a.options.client.Finalize(a.finalize, csr)
	if err != nil {
		return nil, err
	}
	// 获取证书
	cert, err := a.options.client.GetCertificate(order.Certificate)
	if err != nil {
		return nil, err
	}
	return cert, nil
}
