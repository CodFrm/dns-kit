package cdn_svc

import (
	"context"
	"sync"
	"time"

	"github.com/codfrm/cago/database/db"
	"github.com/codfrm/cago/pkg/consts"
	"github.com/codfrm/cago/pkg/i18n"
	"github.com/codfrm/cago/pkg/iam/audit"
	"github.com/codfrm/cago/pkg/logger"
	"github.com/codfrm/cago/pkg/utils/httputils"
	"github.com/codfrm/dns-kit/internal/model/entity/cdn_entity"
	"github.com/codfrm/dns-kit/internal/model/entity/provider_entity"
	"github.com/codfrm/dns-kit/internal/pkg/code"
	"github.com/codfrm/dns-kit/internal/repository/cdn_repo"
	"github.com/codfrm/dns-kit/internal/repository/provider_repo"
	"github.com/codfrm/dns-kit/pkg/platform"
	"go.uber.org/zap"
	"gorm.io/gorm"

	api "github.com/codfrm/dns-kit/internal/api/cdn"
)

type CdnSvc interface {
	// Query 查询cdn
	Query(ctx context.Context, req *api.QueryRequest) (*api.QueryResponse, error)
	// List 获取纳管的cdn列表
	List(ctx context.Context, req *api.ListRequest) (*api.ListResponse, error)
	// Add 添加cdn进入纳管
	Add(ctx context.Context, req *api.AddRequest) (*api.AddResponse, error)
	// Delete 删除cdn
	Delete(ctx context.Context, req *api.DeleteRequest) (*api.DeleteResponse, error)
}

type cdnSvc struct {
	sync.Mutex
}

var defaultCdn = &cdnSvc{}

func Cdn() CdnSvc {
	return defaultCdn
}

// Query 查询cdn
func (c *cdnSvc) Query(ctx context.Context, req *api.QueryRequest) (*api.QueryResponse, error) {
	provider, _, err := provider_repo.Provider().FindPage(ctx, httputils.PageRequest{}, func(filter *provider_repo.FindPageFilters) {
		filter.Platform = []provider_entity.Platform{
			provider_entity.PlatformTencent,
			provider_entity.PlatformQiniu,
		}
	})
	if err != nil {
		return nil, err
	}
	// 获取纳管cdn
	cdnList, _, err := cdn_repo.Cdn().FindPage(ctx, httputils.PageRequest{
		Size: 100,
	})
	if err != nil {
		return nil, err
	}
	cdnMap := make(map[string]bool)
	for _, v := range cdnList {
		cdnMap[v.CdnID] = true
	}
	items := make([]*api.QueryItem, 0)
	for _, p := range provider {
		manager, err := p.CDNManger(ctx)
		if err != nil {
			logger.Ctx(ctx).Warn("cdn manager error", zap.Error(err))
			continue
		}
		list, err := manager.GetCDNList(ctx)
		if err != nil {
			logger.Ctx(ctx).Warn("cdn manager error", zap.Error(err))
			continue
		}
		for _, d := range list {
			items = append(items, &api.QueryItem{
				ProviderID:   p.ID,
				ProviderName: p.Name,
				ID:           d.ID,
				Domain:       d.Domain,
				IsManaged:    cdnMap[d.ID],
			})
		}
	}
	return &api.QueryResponse{
		Items: items,
	}, nil
}

// List 获取纳管的cdn列表
func (c *cdnSvc) List(ctx context.Context, req *api.ListRequest) (*api.ListResponse, error) {
	list, total, err := cdn_repo.Cdn().FindPage(ctx, req.PageRequest)
	if err != nil {
		return nil, err
	}
	items := make([]*api.Item, 0)
	for _, v := range list {
		provider, err := provider_repo.Provider().Find(ctx, v.ProviderID)
		if err != nil {
			return nil, err
		}
		items = append(items, &api.Item{
			ID:           v.ID,
			ProviderName: provider.Name,
			Domain:       v.Domain,
			Createtime:   v.Createtime,
		})
	}
	return &api.ListResponse{
		PageResponse: httputils.PageResponse[*api.Item]{
			List:  items,
			Total: total,
		},
	}, nil
}

// Add 添加cdn进入纳管
func (c *cdnSvc) Add(ctx context.Context, req *api.AddRequest) (*api.AddResponse, error) {
	c.Lock()
	defer c.Unlock()
	provider, err := provider_repo.Provider().Find(ctx, req.ProviderID)
	if err != nil {
		return nil, err
	}
	if err := provider.Check(ctx); err != nil {
		return nil, err
	}
	manager, err := provider.CDNManger(ctx)
	if err != nil {
		return nil, err
	}
	detail, err := manager.GetCDNDetail(ctx, &platform.CDNItem{
		ID:     req.ID,
		Domain: req.Domain,
	})
	if err != nil {
		return nil, err
	}
	if detail == nil {
		return nil, i18n.NewError(ctx, code.CDNNotFound)
	}
	// 纳管
	cdn := &cdn_entity.Cdn{
		ProviderID: req.ProviderID,
		CdnID:      req.ID,
		Domain:     req.Domain,
		Status:     consts.ACTIVE,
		Createtime: time.Now().Unix(),
	}
	if err := cdn_repo.Cdn().Create(ctx, cdn); err != nil {
		return nil, err
	}
	_ = audit.Ctx(ctx).Record("create", zap.Int64("id", cdn.ID), zap.String("name", cdn.Domain))
	return nil, nil
}

// Delete 删除cdn
func (c *cdnSvc) Delete(ctx context.Context, req *api.DeleteRequest) (*api.DeleteResponse, error) {
	c.Lock()
	defer c.Unlock()
	cdn, err := cdn_repo.Cdn().Find(ctx, req.ID)
	if err != nil {
		return nil, err
	}
	if err := cdn.Check(ctx); err != nil {
		return nil, err
	}
	err = db.Ctx(ctx).Transaction(func(tx *gorm.DB) error {
		ctx = db.WithContextDB(ctx, tx)
		if err := cdn_repo.Cdn().Delete(ctx, cdn.ID); err != nil {
			return err
		}
		// 删除相关资源
		return nil
	})
	if err != nil {
		return nil, err
	}
	_ = audit.Ctx(ctx).Record("delete", zap.Int64("id", cdn.ID), zap.String("name", cdn.Domain))
	return nil, nil
}
