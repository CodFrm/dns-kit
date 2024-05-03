package provider_svc

import (
	"context"
	"encoding/json"
	"github.com/codfrm/dns-kit/internal/api/cdn"
	"github.com/codfrm/dns-kit/internal/api/domain"
	"github.com/codfrm/dns-kit/internal/repository/cdn_repo"
	"github.com/codfrm/dns-kit/internal/repository/domain_repo"
	"github.com/codfrm/dns-kit/internal/service/cdn_svc"
	"github.com/codfrm/dns-kit/internal/service/domain_svc"
	"sync"
	"time"

	"github.com/codfrm/cago/database/db"
	"github.com/codfrm/cago/pkg/iam/audit"
	"github.com/codfrm/cago/pkg/utils/httputils"
	"go.uber.org/zap"
	"gorm.io/gorm"

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
	sync.Mutex
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
	p.Lock()
	defer p.Unlock()
	secretData, err := json.Marshal(req.Secret)
	if err != nil {
		return nil, err
	}
	provider2 := &provider_entity.Provider{
		Name: req.Name,
		//UserID:     user.ID,
		Secret:     string(secretData),
		Platform:   req.Platform,
		Status:     consts.ACTIVE,
		Createtime: time.Now().Unix(),
		Updatetime: time.Now().Unix(),
	}
	manager, err := provider2.DomainManager(ctx)
	if err != nil {
		return nil, err
	}
	user, err := manager.UserDetails(ctx)
	if err != nil {
		return nil, i18n.NewError(ctx, code.ProviderSecretError)
	}
	provider2.UserID = user.ID
	// 判断重复添加
	provider, err := provider_repo.Provider().FindByProviderUserId(ctx, provider2.UserID)
	if err != nil {
		return nil, err
	}
	if len(provider) > 0 {
		return nil, i18n.NewError(ctx, code.ProviderExist)
	}

	if err := provider_repo.Provider().Create(ctx, provider2); err != nil {
		return nil, err
	}
	_ = audit.Ctx(ctx).Record("create", zap.String("name", req.Name), zap.String("platform", string(req.Platform)))
	return &api.CreateProviderResponse{}, nil
}

// UpdateProvider 更新供应商
func (p *providerSvc) UpdateProvider(ctx context.Context, req *api.UpdateProviderRequest) (*api.UpdateProviderResponse, error) {
	p.Lock()
	defer p.Unlock()
	provider, err := provider_repo.Provider().Find(ctx, req.ID)
	if err != nil {
		return nil, err
	}
	if err := provider.Check(ctx); err != nil {
		return nil, err
	}
	provider.Name = req.Name
	if len(req.Secret) > 0 {
		manager, err := provider.DomainManager(ctx)
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
	_ = audit.Ctx(ctx).Record("update", zap.String("name", provider.Name))
	return nil, nil
}

// DeleteProvider 删除供应商
func (p *providerSvc) DeleteProvider(ctx context.Context, req *api.DeleteProviderRequest) (*api.DeleteProviderResponse, error) {
	p.Lock()
	defer p.Unlock()
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
		// 删除相关的域名和cdn
		domainList, err := domain_repo.Domain().FindByProviderId(ctx, req.ID)
		if err != nil {
			return err
		}
		for _, v := range domainList {
			if _, err := domain_svc.Domain().Delete(ctx, &domain.DeleteRequest{ID: v.ID}); err != nil {
				return err
			}
		}
		cdnList, err := cdn_repo.Cdn().FindByProviderId(ctx, req.ID)
		if err != nil {
			return err
		}
		for _, v := range cdnList {
			if _, err := cdn_svc.Cdn().Delete(ctx, &cdn.DeleteRequest{ID: v.ID}); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	_ = audit.Ctx(ctx).Record("delete", zap.String("name", provider.Name))
	return nil, nil
}
