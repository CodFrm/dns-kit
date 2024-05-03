package cert_entity

import (
	"context"
	"database/sql/driver"
	"encoding/json"

	"github.com/codfrm/cago/pkg/i18n"
	"github.com/codfrm/dns-kit/internal/pkg/code"
)

type CertStatus int32

const (
	CertStatusActive    = 1 + iota // 激活状态
	CertStatusDel                  // 删除状态
	CertStatusExpire               // 过期状态
	CertStatusApply                // 申请中
	CertStatusApplyFail            // 申请失败
)

type GORMStringSlice []string

func (d *GORMStringSlice) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	return json.Unmarshal(value.([]byte), d)
}

func (d GORMStringSlice) Value() (driver.Value, error) {
	if d == nil {
		return nil, nil
	}
	return json.Marshal(d)
}

type Cert struct {
	ID                 int64           `gorm:"column:id;not null;primary_key"`
	Email              string          `gorm:"column:email;type:varchar(255);not null"`          // 邮箱
	Domains            GORMStringSlice `gorm:"column:domains;type:varchar(255)"`                 // 域名列表
	Certificate        string          `gorm:"column:certificate;type:text"`                     // 证书
	CertificateRequest string          `gorm:"column:certificate_request;type:text"`             // 证书请求
	PrivateKey         string          `gorm:"column:private_key;type:text"`                     // 私钥
	Status             CertStatus      `gorm:"column:status;type:tinyint(4);default:0;not null"` // 状态
	Createtime         int64           `gorm:"column:createtime;type:bigint(20)"`
	Updatetime         int64           `gorm:"column:updatetime;type:bigint(20)"`
	Expiretime         int64           `gorm:"column:expiretime;type:bigint(20)"`
}

func (c *Cert) Check(ctx context.Context) error {
	if c == nil {
		return i18n.NewError(ctx, code.CertNotFound)
	}
	return nil
}
