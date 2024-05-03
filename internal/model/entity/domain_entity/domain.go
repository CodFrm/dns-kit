package domain_entity

import (
	"context"

	"github.com/codfrm/cago/pkg/i18n"
	"github.com/codfrm/dns-kit/internal/pkg/code"
	"github.com/codfrm/dns-kit/internal/repository/provider_repo"
	"github.com/codfrm/dns-kit/pkg/platform"
)

type Domain struct {
	ID         int64  `gorm:"column:id;not null;primary_key"`
	ProviderID int64  `gorm:"column:provider_id;type:bigint(20);not null"` // 供应商id
	DomainID   string `gorm:"column:domain_id;type:varchar(128);not null"` // 域名id
	Domain     string `gorm:"column:domain;type:varchar(128);not null"`    // 域名
	Status     int32  `gorm:"column:status;type:tinyint(4);default:0;not null"`
	Createtime int64  `gorm:"column:createtime;type:bigint(20)"`
	Updatetime int64  `gorm:"column:updatetime;type:bigint(20)"`
}

func (d *Domain) Check(ctx context.Context) error {
	if d == nil {
		return i18n.NewError(ctx, code.DomainNotFound)
	}
	return nil
}

func (d *Domain) DnsManager(ctx context.Context) (platform.DNSManager, error) {
	provider, err := provider_repo.Provider().Find(ctx, d.ProviderID)
	if err != nil {
		return nil, err
	}
	if err := provider.Check(ctx); err != nil {
		return nil, err
	}
	manager, err := provider.DomainManager(ctx)
	if err != nil {
		return nil, err
	}
	dnsManager, err := manager.BuildDNSManager(ctx, &platform.Domain{
		ID:     d.DomainID,
		Domain: d.Domain,
	})
	if err != nil {
		return nil, err
	}
	return dnsManager, nil
}
