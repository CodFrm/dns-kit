package acme_entity

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"math/big"

	"github.com/codfrm/dns-kit/pkg/acme"
)

// Acme acme账号
type Acme struct {
	ID         int64  `gorm:"column:id;not null;primary_key"`
	Email      string `gorm:"column:email;type:varchar(255);index:email,unique;not null"` // 邮箱
	Kid        string `gorm:"column:kid;type:varchar(255);not null"`                      // kid
	PrivateKey string `gorm:"column:private_key;type:varchar(255);not null"`              // 私钥
	Status     int8   `gorm:"column:status;type:tinyint(4);not null"`                     // 状态
	Createtime int64  `gorm:"column:createtime;type:bigint(20);not null"`                 // 创建时间
}

type PrivateKeyJson struct {
	PX *big.Int
	PY *big.Int
	D  *big.Int
}

func (a *Acme) NewACME() (*acme.Acme, error) {
	var client *acme.Client
	block, _ := pem.Decode([]byte(a.PrivateKey))
	if block == nil {
		return nil, errors.New("failed to parse private key PEM")
	}
	privateKey, err := x509.ParseECPrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	client, err = acme.NewClient(acme.WithKid(a.Kid), acme.WithPrivateKey(privateKey))
	if err != nil {
		return nil, err
	}
	return acme.NewAcme(a.Email, acme.WithClient(client))
}

func (a *Acme) SavePrivateKey(privateKey *ecdsa.PrivateKey) error {
	data, err := x509.MarshalECPrivateKey(privateKey)
	if err != nil {
		return err
	}
	buf := bytes.NewBuffer(nil)
	if err := pem.Encode(buf, &pem.Block{Type: "PRIVATE KEY", Bytes: data}); err != nil {
		return err
	}
	a.PrivateKey = buf.String()
	return nil
}
