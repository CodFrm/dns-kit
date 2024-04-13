package dns_provider_repo

import (
	"context"

	"github.com/codfrm/cago/database/db"
	"github.com/codfrm/cago/pkg/consts"
	"github.com/codfrm/cago/pkg/utils/httputils"
	"github.com/codfrm/dns-kit/internal/model/entity/dns_provider_entity"
)

type DnsProviderRepo interface {
	Find(ctx context.Context, id int64) (*dns_provider_entity.DnsProvider, error)
	FindPage(ctx context.Context, page httputils.PageRequest) ([]*dns_provider_entity.DnsProvider, int64, error)
	Create(ctx context.Context, dnsProvider *dns_provider_entity.DnsProvider) error
	Update(ctx context.Context, dnsProvider *dns_provider_entity.DnsProvider) error
	Delete(ctx context.Context, id int64) error

	FindByProviderUserId(ctx context.Context, userId string) (*dns_provider_entity.DnsProvider, error)
}

var defaultDnsProvider DnsProviderRepo

func DnsProvider() DnsProviderRepo {
	return defaultDnsProvider
}

func RegisterDnsProvider(i DnsProviderRepo) {
	defaultDnsProvider = i
}

type dnsProviderRepo struct {
}

func NewDnsProvider() DnsProviderRepo {
	return &dnsProviderRepo{}
}

func (u *dnsProviderRepo) Find(ctx context.Context, id int64) (*dns_provider_entity.DnsProvider, error) {
	ret := &dns_provider_entity.DnsProvider{}
	if err := db.Ctx(ctx).Where("id=? and status=?", id, consts.ACTIVE).First(ret).Error; err != nil {
		if db.RecordNotFound(err) {
			return nil, nil
		}
		return nil, err
	}
	return ret, nil
}

func (u *dnsProviderRepo) Create(ctx context.Context, dnsProvider *dns_provider_entity.DnsProvider) error {
	return db.Ctx(ctx).Create(dnsProvider).Error
}

func (u *dnsProviderRepo) Update(ctx context.Context, dnsProvider *dns_provider_entity.DnsProvider) error {
	return db.Ctx(ctx).Updates(dnsProvider).Error
}

func (u *dnsProviderRepo) Delete(ctx context.Context, id int64) error {
	return db.Ctx(ctx).Model(&dns_provider_entity.DnsProvider{}).Where("id=?", id).Update("status", consts.DELETE).Error
}

func (u *dnsProviderRepo) FindPage(ctx context.Context, page httputils.PageRequest) ([]*dns_provider_entity.DnsProvider, int64, error) {
	var list []*dns_provider_entity.DnsProvider
	var count int64
	find := db.Ctx(ctx).Model(&dns_provider_entity.DnsProvider{}).Where("status=?", consts.ACTIVE)
	if err := find.Count(&count).Error; err != nil {
		return nil, 0, err
	}
	if err := find.Order("createtime desc").Offset(page.GetOffset()).Limit(page.GetLimit()).Find(&list).Error; err != nil {
		return nil, 0, err
	}
	return list, count, nil
}

func (u *dnsProviderRepo) FindByProviderUserId(ctx context.Context, userId string) (*dns_provider_entity.DnsProvider, error) {
	ret := &dns_provider_entity.DnsProvider{}
	if err := db.Ctx(ctx).Where("user_id=? and status=?", userId, consts.ACTIVE).First(ret).Error; err != nil {
		if db.RecordNotFound(err) {
			return nil, nil
		}
		return nil, err
	}
	return ret, nil
}
