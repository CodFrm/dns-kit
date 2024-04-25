package provider_repo

import (
	"context"

	"github.com/codfrm/cago/database/db"
	"github.com/codfrm/cago/pkg/consts"
	"github.com/codfrm/cago/pkg/utils/httputils"
	"github.com/codfrm/dns-kit/internal/model/entity/provider_entity"
)

type ProviderRepo interface {
	Find(ctx context.Context, id int64) (*provider_entity.Provider, error)
	FindPage(ctx context.Context, page httputils.PageRequest) ([]*provider_entity.Provider, int64, error)
	Create(ctx context.Context, dnsProvider *provider_entity.Provider) error
	Update(ctx context.Context, dnsProvider *provider_entity.Provider) error
	Delete(ctx context.Context, id int64) error

	FindByProviderUserId(ctx context.Context, userId string) (*provider_entity.Provider, error)
}

var defaultDnsProvider ProviderRepo

func DnsProvider() ProviderRepo {
	return defaultDnsProvider
}

func RegisterDnsProvider(i ProviderRepo) {
	defaultDnsProvider = i
}

type dnsProviderRepo struct {
}

func NewDnsProvider() ProviderRepo {
	return &dnsProviderRepo{}
}

func (u *dnsProviderRepo) Find(ctx context.Context, id int64) (*provider_entity.Provider, error) {
	ret := &provider_entity.Provider{}
	if err := db.Ctx(ctx).Where("id=? and status=?", id, consts.ACTIVE).First(ret).Error; err != nil {
		if db.RecordNotFound(err) {
			return nil, nil
		}
		return nil, err
	}
	return ret, nil
}

func (u *dnsProviderRepo) Create(ctx context.Context, dnsProvider *provider_entity.Provider) error {
	return db.Ctx(ctx).Create(dnsProvider).Error
}

func (u *dnsProviderRepo) Update(ctx context.Context, dnsProvider *provider_entity.Provider) error {
	return db.Ctx(ctx).Updates(dnsProvider).Error
}

func (u *dnsProviderRepo) Delete(ctx context.Context, id int64) error {
	return db.Ctx(ctx).Model(&provider_entity.Provider{}).Where("id=?", id).Update("status", consts.DELETE).Error
}

func (u *dnsProviderRepo) FindPage(ctx context.Context, page httputils.PageRequest) ([]*provider_entity.Provider, int64, error) {
	var list []*provider_entity.Provider
	var count int64
	find := db.Ctx(ctx).Model(&provider_entity.Provider{}).Where("status=?", consts.ACTIVE)
	if err := find.Count(&count).Error; err != nil {
		return nil, 0, err
	}
	if err := find.Order("createtime desc").Offset(page.GetOffset()).Limit(page.GetLimit()).Find(&list).Error; err != nil {
		return nil, 0, err
	}
	return list, count, nil
}

func (u *dnsProviderRepo) FindByProviderUserId(ctx context.Context, userId string) (*provider_entity.Provider, error) {
	ret := &provider_entity.Provider{}
	if err := db.Ctx(ctx).Where("user_id=? and status=?", userId, consts.ACTIVE).First(ret).Error; err != nil {
		if db.RecordNotFound(err) {
			return nil, nil
		}
		return nil, err
	}
	return ret, nil
}
