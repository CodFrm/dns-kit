package cert_hosting_entity

import (
	"context"
	"encoding/json"

	"github.com/codfrm/cago/pkg/i18n"
	"github.com/codfrm/dns-kit/internal/pkg/code"
)

type CertHostingStatus int

const (
	CertHostingStatusActive CertHostingStatus = iota + 1
	CertHostingStatusDelete
	CertHostingStatusDeploy     // 部署中
	CertHostingStatusDeployFail // 部署失败
	CertHostingStatusFail       // 托管失败(不可重试)
)

// CertHostingType 证书托管类型 1:CDN 2:供应商
type CertHostingType int

const (
	CertHostingTypeCDN = iota + 1
	CertHostingTypeProvider
)

// CertHosting 证书托管
type CertHosting struct {
	ID         int64           `gorm:"column:id;not null;primary_key"`
	Email      string          `gorm:"column:email;type:varchar(255);not null"`        // 邮箱
	Type       CertHostingType `gorm:"column:type;type:tinyint(4);default:1;not null"` // 类型
	CdnID      int64           `gorm:"column:cdn_id;type:bigint(20);"`                 // CDN ID，当类型为 CDN 时，该字段不能为空
	ProviderID int64           `gorm:"column:provider_id;type:bigint(20);"`            // 供应商 ID，当类型为供应商时，该字段不能为空
	Config     string          `gorm:"column:config;type:varchar(255);"`               // 域名，当类型为供应商时，该字段不能为空
	configMap  map[string]string
	CertID     int64             `gorm:"column:cert_id;type:bigint(20);default:0"`         // 关联的证书 ID
	Status     CertHostingStatus `gorm:"column:status;type:tinyint(4);default:0;not null"` // 状态
	Createtime int64             `gorm:"column:createtime;type:bigint(20)"`
	Updatetime int64             `gorm:"column:updatetime;type:bigint(20)"`
}

func (h *CertHosting) Check(ctx context.Context) error {
	if h == nil {
		return i18n.NewError(ctx, code.CertHostingNotFound)
	}
	return nil
}

func (h *CertHosting) SetConfigMap(config map[string]string) error {
	h.configMap = config
	data, err := json.Marshal(config)
	if err != nil {
		return err
	}
	h.Config = string(data)
	return nil
}

func (h *CertHosting) ConfigMap() map[string]string {
	if h.configMap != nil {
		return h.configMap
	}
	ret := make(map[string]string)
	if h.Config != "" {
		_ = json.Unmarshal([]byte(h.Config), &ret)
	}
	h.configMap = ret
	return ret
}
