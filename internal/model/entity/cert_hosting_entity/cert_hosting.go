package cert_hosting_entity

import (
	"context"

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

// CertHosting 证书托管
type CertHosting struct {
	ID         int64             `gorm:"column:id;type:bigint(20);not null;primary_key"`
	Email      string            `gorm:"column:email;type:varchar(255);not null"`          // 邮箱
	CdnID      int64             `gorm:"column:cdn_id;type:bigint(20);not null"`           // CDN ID
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
