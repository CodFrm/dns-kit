package cdn_entity

import (
	"context"

	"github.com/codfrm/cago/pkg/i18n"
	"github.com/codfrm/dns-kit/internal/pkg/code"
	"github.com/codfrm/dns-kit/internal/repository/provider_repo"
	"github.com/codfrm/dns-kit/pkg/platform"
)

type Cdn struct {
	ID         int64  `gorm:"column:id;type:bigint(20);not null;primary_key"`
	ProviderID int64  `gorm:"column:provider_id;type:bigint(20);not null"`      // 供应商id
	CdnID      string `gorm:"column:cdn_id;type:varchar(128);not null"`         // cdn id
	Domain     string `gorm:"column:domain;type:varchar(128);not null"`         // 域名
	Status     int32  `gorm:"column:status;type:tinyint(4);default:0;not null"` // 状态
	Createtime int64  `gorm:"column:createtime;type:bigint(20)"`
	Updatetime int64  `gorm:"column:updatetime;type:bigint(20)"`
}

func (c *Cdn) Check(ctx context.Context) error {
	if c == nil {
		return i18n.NewError(ctx, code.CDNNotFound)
	}
	return nil
}

func (c *Cdn) CDNManger(ctx context.Context) (platform.CDNManager, error) {
	provider, err := provider_repo.Provider().Find(ctx, c.ProviderID)
	if err != nil {
		return nil, err
	}
	if err := provider.Check(ctx); err != nil {
		return nil, err
	}
	return provider.CDNManger(ctx)
}
