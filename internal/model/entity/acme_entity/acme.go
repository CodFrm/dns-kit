package acme_entity

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"encoding/json"
	"github.com/codfrm/dns-kit/pkg/acme"
	"math/big"
)

// Acme acme账号
type Acme struct {
	ID         int64  `gorm:"column:id;type:bigint(20);not null;primary_key"`
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
	privateKeyJson := &PrivateKeyJson{}
	err := json.Unmarshal([]byte(a.PrivateKey), privateKeyJson)
	if err != nil {
		return nil, err
	}
	client, err = acme.NewClient(acme.WithKid(a.Kid), acme.WithPrivateKey(&ecdsa.PrivateKey{
		PublicKey: ecdsa.PublicKey{
			Curve: elliptic.P256(),
			X:     privateKeyJson.PX,
			Y:     privateKeyJson.PY,
		},
		D: privateKeyJson.D,
	}))
	if err != nil {
		return nil, err
	}
	return acme.NewAcme(a.Email, acme.WithClient(client))
}

func (a *Acme) SavePrivateKey(privateKey *ecdsa.PrivateKey) error {
	privateKeyData, err := json.Marshal(&PrivateKeyJson{
		PX: privateKey.PublicKey.X,
		PY: privateKey.PublicKey.Y,
		D:  privateKey.D,
	})
	if err != nil {
		return err
	}
	a.PrivateKey = string(privateKeyData)
	return nil
}
