package provider_svc

import (
	"context"
	"encoding/json"
	"github.com/codfrm/cago/database/db"
	"github.com/codfrm/cago/pkg/iam/audit"
	"github.com/codfrm/cago/pkg/utils/httputils"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"time"

	"github.com/codfrm/cago/pkg/consts"
	"github.com/codfrm/cago/pkg/i18n"
	api "github.com/codfrm/dns-kit/internal/api/provider"
	"github.com/codfrm/dns-kit/internal/model/entity/provider_entity"
	"github.com/codfrm/dns-kit/internal/pkg/code"
	"github.com/codfrm/dns-kit/internal/repository/provider_repo"
)

type ProviderSvc interface {
	// ListProvider 获取供应商列表
	ListProvider(ctx context.Context, req *api.ListProviderRequest) (*api.ListProviderResponse, error)
	// CreateProvider 创建供应商
	CreateProvider(ctx context.Context, req *api.CreateProviderRequest) (*api.CreateProviderResponse, error)
	// UpdateProvider 更新供应商
	UpdateProvider(ctx context.Context, req *api.UpdateProviderRequest) (*api.UpdateProviderResponse, error)
	// DeleteProvider 删除供应商
	DeleteProvider(ctx context.Context, req *api.DeleteProviderRequest) (*api.DeleteProviderResponse, error)
}

type providerSvc struct {
}

var defaultProvider = &providerSvc{}

func Provider() ProviderSvc {
	return defaultProvider
}

// ListProvider 获取供应商列表
func (p *providerSvc) ListProvider(ctx context.Context, req *api.ListProviderRequest) (*api.ListProviderResponse, error) {
	list, total, err := provider_repo.Provider().FindPage(ctx, req.PageRequest)
	if err != nil {
		return nil, err
	}
	ret := &api.ListProviderResponse{
		PageResponse: httputils.PageResponse[*api.Item]{
			List:  make([]*api.Item, 0),
			Total: total,
		},
	}
	for _, v := range list {
		ret.List = append(ret.List, &api.Item{
			ID:       v.ID,
			Name:     v.Name,
			Platform: v.Platform,
		})
	}
	return ret, nil
}

// CreateProvider 创建供应商
func (p *providerSvc) CreateProvider(ctx context.Context, req *api.CreateProviderRequest) (*api.CreateProviderResponse, error) {
	manager, err := NewProvider(ctx, req.Platform, req.Secret)
	if err != nil {
		return nil, err
	}
	user, err := manager.UserDetails(ctx)
	if err != nil {
		return nil, i18n.NewError(ctx, code.ProviderSecretError)
	}
	// 判断重复添加
	provider, err := provider_repo.Provider().FindByProviderUserId(ctx, user.ID)
	if err != nil {
		return nil, err
	}
	if len(provider) > 0 {
		return nil, i18n.NewError(ctx, code.ProviderExist)
	}
	secretData, err := json.Marshal(req.Secret)
	if err != nil {
		return nil, err
	}
	provider2 := &provider_entity.Provider{
		Name:       req.Name,
		UserID:     user.ID,
		Secret:     string(secretData),
		Platform:   req.Platform,
		Status:     consts.ACTIVE,
		Createtime: time.Now().Unix(),
		Updatetime: time.Now().Unix(),
	}
	if err := provider_repo.Provider().Create(ctx, provider2); err != nil {
		return nil, err
	}
	_ = audit.Ctx(ctx).Record(ctx, "create", zap.String("name", req.Name), zap.String("platform", string(req.Platform)))
	return &api.CreateProviderResponse{}, nil
}

// UpdateProvider 更新供应商
func (p *providerSvc) UpdateProvider(ctx context.Context, req *api.UpdateProviderRequest) (*api.UpdateProviderResponse, error) {
	provider, err := provider_repo.Provider().Find(ctx, req.ID)
	if err != nil {
		return nil, err
	}
	if err := provider.Check(ctx); err != nil {
		return nil, err
	}
	provider.Name = req.Name
	if len(req.Secret) > 0 {
		manager, err := NewProvider(ctx, provider.Platform, req.Secret)
		if err != nil {
			return nil, err
		}
		user, err := manager.UserDetails(ctx)
		if err != nil {
			return nil, i18n.NewError(ctx, code.ProviderSecretError)
		}
		// 判断重复添加
		provides, err := provider_repo.Provider().FindByProviderUserId(ctx, user.ID)
		if err != nil {
			return nil, err
		}
		for _, v := range provides {
			if v.ID != req.ID {
				return nil, i18n.NewError(ctx, code.ProviderExist)
			}
		}
		secretData, err := json.Marshal(req.Secret)
		if err != nil {
			return nil, err
		}
		provider.UserID = user.ID
		provider.Secret = string(secretData)
	}
	provider.Updatetime = time.Now().Unix()
	if err := provider_repo.Provider().Update(ctx, provider); err != nil {
		return nil, err
	}
	_ = audit.Ctx(ctx).Record(ctx, "update", zap.String("name", provider.Name))
	return nil, nil
}

// DeleteProvider 删除供应商
func (p *providerSvc) DeleteProvider(ctx context.Context, req *api.DeleteProviderRequest) (*api.DeleteProviderResponse, error) {
	provider, err := provider_repo.Provider().Find(ctx, req.ID)
	if err != nil {
		return nil, err
	}
	if err := provider.Check(ctx); err != nil {
		return nil, err
	}
	err = db.Ctx(ctx).Transaction(func(tx *gorm.DB) error {
		ctx = db.WithContextDB(ctx, tx)
		if err := provider_repo.Provider().Delete(ctx, req.ID); err != nil {
			return err
		}
		// 删除相关资源
		return nil
	})
	if err != nil {
		return nil, err
	}
	_ = audit.Ctx(ctx).Record(ctx, "delete", zap.String("name", provider.Name))
	return nil, nil
}