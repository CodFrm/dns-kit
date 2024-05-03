package domain_svc

import (
	"context"
	"github.com/codfrm/dns-kit/internal/model/entity/provider_entity"
	"sync"
	"time"

	"github.com/codfrm/cago/database/db"
	"github.com/codfrm/cago/pkg/consts"
	"github.com/codfrm/cago/pkg/i18n"
	"github.com/codfrm/cago/pkg/iam/audit"
	"github.com/codfrm/cago/pkg/logger"
	"github.com/codfrm/cago/pkg/utils/httputils"
	api "github.com/codfrm/dns-kit/internal/api/domain"
	"github.com/codfrm/dns-kit/internal/model/entity/domain_entity"
	"github.com/codfrm/dns-kit/internal/pkg/code"
	"github.com/codfrm/dns-kit/internal/repository/domain_repo"
	"github.com/codfrm/dns-kit/internal/repository/provider_repo"
	"github.com/codfrm/dns-kit/pkg/platform"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type DomainSvc interface {
	// List 获取域名列表
	List(ctx context.Context, req *api.ListRequest) (*api.ListResponse, error)
	// Query 查询域名列表
	Query(ctx context.Context, req *api.QueryRequest) (*api.QueryResponse, error)
	// Add 纳管域名
	Add(ctx context.Context, req *api.AddRequest) (*api.AddResponse, error)
	// Delete 删除域名
	Delete(ctx context.Context, req *api.DeleteRequest) (*api.DeleteResponse, error)
}

type domainSvc struct {
	sync.Mutex
}

var defaultDomain = &domainSvc{}

func Domain() DomainSvc {
	return defaultDomain
}

// List 获取域名列表
func (d *domainSvc) List(ctx context.Context, req *api.ListRequest) (*api.ListResponse, error) {
	list, total, err := domain_repo.Domain().FindPage(ctx, req.PageRequest)
	if err != nil {
		return nil, err
	}
	ret := &api.ListResponse{
		PageResponse: httputils.PageResponse[*api.Item]{
			List:  make([]*api.Item, 0),
			Total: total,
		},
	}
	for _, v := range list {
		provider, err := provider_repo.Provider().Find(ctx, v.ProviderID)
		if err != nil {
			return nil, err
		}
		domain := &api.Item{
			ID:           v.ID,
			Domain:       v.Domain,
			ProviderName: provider.Name,
		}
		ret.List = append(ret.List, domain)
	}
	return ret, nil
}

// Query 查询域名列表
func (d *domainSvc) Query(ctx context.Context, req *api.QueryRequest) (*api.QueryResponse, error) {
	// 获取供应商信息
	provides, _, err := provider_repo.Provider().FindPage(ctx, httputils.PageRequest{
		Size: 100,
	}, func(filter *provider_repo.FindPageFilters) {
		filter.Platform = []provider_entity.Platform{provider_entity.PlatformTencent, provider_entity.PlatformCloudflare}
	})
	if err != nil {
		return nil, err
	}
	ret := &api.QueryResponse{
		Items: make([]*api.QueryItem, 0),
	}
	// 获取纳管域名
	domains, _, err := domain_repo.Domain().FindPage(ctx, httputils.PageRequest{
		Size: 100,
	})
	if err != nil {
		return nil, err
	}
	domainMap := make(map[string]bool)
	for _, v := range domains {
		domainMap[v.DomainID] = true
	}
	// 获取域名信息
	for _, v := range provides {
		manager, err := v.DomainManager(ctx)
		if err != nil {
			logger.Ctx(ctx).Error("provider_svc.NewProvider", zap.Error(err))
			continue
		}
		// 获取域名列表
		domainList, err := manager.GetDomainList(ctx)
		if err != nil {
			logger.Ctx(ctx).Error("manager.GetDomainList", zap.Error(err))
			continue
		}
		for _, domain := range domainList {
			ret.Items = append(ret.Items, &api.QueryItem{
				ProviderID:   v.ID,
				ProviderName: v.Name,
				DomainID:     domain.ID,
				Domain:       domain.Domain,
				IsManaged:    domainMap[domain.ID],
			})
		}
	}
	return ret, nil
}

// Add 纳管域名
func (d *domainSvc) Add(ctx context.Context, req *api.AddRequest) (*api.AddResponse, error) {
	d.Lock()
	defer d.Unlock()
	// 判断是否已经纳管
	domain, err := domain_repo.Domain().FindByDomainID(ctx, req.DomainID)
	if err != nil {
		return nil, err
	}
	if domain != nil {
		return nil, i18n.NewError(ctx, code.DomainIsManaged)
	}
	// 获取供应商信息
	provider, err := provider_repo.Provider().Find(ctx, req.ProviderID)
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
		ID:     req.DomainID,
		Domain: req.Domain,
	})
	if err != nil {
		return nil, err
	}
	_, err = dnsManager.GetRecordList(ctx)
	if err != nil {
		return nil, err
	}
	// 加入纳管
	domain = &domain_entity.Domain{
		ProviderID: provider.ID,
		DomainID:   req.DomainID,
		Domain:     req.Domain,
		Status:     consts.ACTIVE,
		Createtime: time.Now().Unix(),
	}
	if err := domain_repo.Domain().Create(ctx, domain); err != nil {
		return nil, err
	}
	_ = audit.Ctx(ctx).Record("create", zap.String("domain", req.Domain))
	return nil, nil
}

// Delete 删除域名
func (d *domainSvc) Delete(ctx context.Context, req *api.DeleteRequest) (*api.DeleteResponse, error) {
	d.Lock()
	defer d.Unlock()
	domain, err := domain_repo.Domain().Find(ctx, req.ID)
	if err != nil {
		return nil, err
	}
	if err := domain.Check(ctx); err != nil {
		return nil, err
	}
	// 删除相关资源
	err = db.Ctx(ctx).Transaction(func(tx *gorm.DB) error {
		ctx = db.WithContextDB(ctx, tx)
		if err := domain_repo.Domain().Delete(ctx, req.ID); err != nil {
			return err
		}
		// 删除相关资源
		return nil
	})
	if err != nil {
		return nil, err
	}
	_ = audit.Ctx(ctx).Record("delete", zap.String("domain", domain.Domain))
	return nil, nil
}
