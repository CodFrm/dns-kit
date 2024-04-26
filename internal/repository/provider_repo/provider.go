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
	Create(ctx context.Context, provider *provider_entity.Provider) error
	Update(ctx context.Context, provider *provider_entity.Provider) error
	Delete(ctx context.Context, id int64) error

	FindByProviderUserId(ctx context.Context, userId string) ([]*provider_entity.Provider, error)
}

var defaultProvider ProviderRepo

func Provider() ProviderRepo {
	return defaultProvider
}

func RegisterProvider(i ProviderRepo) {
	defaultProvider = i
}

type providerRepo struct {
}

func NewProvider() ProviderRepo {
	return &providerRepo{}
}

func (u *providerRepo) Find(ctx context.Context, id int64) (*provider_entity.Provider, error) {
	ret := &provider_entity.Provider{}
	if err := db.Ctx(ctx).Where("id=? and status=?", id, consts.ACTIVE).First(ret).Error; err != nil {
		if db.RecordNotFound(err) {
			return nil, nil
		}
		return nil, err
	}
	return ret, nil
}

func (u *providerRepo) Create(ctx context.Context, provider *provider_entity.Provider) error {
	return db.Ctx(ctx).Create(provider).Error
}

func (u *providerRepo) Update(ctx context.Context, provider *provider_entity.Provider) error {
	return db.Ctx(ctx).Updates(provider).Error
}

func (u *providerRepo) Delete(ctx context.Context, id int64) error {
	return db.Ctx(ctx).Model(&provider_entity.Provider{}).Where("id=?", id).Update("status", consts.DELETE).Error
}

func (u *providerRepo) FindPage(ctx context.Context, page httputils.PageRequest) ([]*provider_entity.Provider, int64, error) {
	var list []*provider_entity.Provider
	var count int64
	find := db.Ctx(ctx).Model(&provider_entity.Provider{}).Where("status=?", consts.ACTIVE)
	if err := find.Count(&count).Error; err != nil {
		return nil, 0, err
	}
	if err := find.Order("createtime desc").
		Offset(page.GetOffset()).Limit(page.GetLimit()).Find(&list).Error; err != nil {
		return nil, 0, err
	}
	return list, count, nil
}

func (u *providerRepo) FindByProviderUserId(ctx context.Context, userId string) ([]*provider_entity.Provider, error) {
	ret := make([]*provider_entity.Provider, 0)
	if err := db.Ctx(ctx).
		Where("user_id=? and status=?", userId, consts.ACTIVE).Find(&ret).Error; err != nil {
		return nil, err
	}
	return ret, nil
}
