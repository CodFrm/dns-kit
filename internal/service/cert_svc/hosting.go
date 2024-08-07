package cert_svc

import (
	"context"
	"sync"
	"time"

	"github.com/codfrm/dns-kit/internal/repository/provider_repo"
	"github.com/codfrm/dns-kit/internal/service/cert_svc/deploy"

	"github.com/codfrm/cago/pkg/i18n"
	"github.com/codfrm/cago/pkg/iam/audit"
	"github.com/codfrm/cago/pkg/logger"
	"github.com/codfrm/cago/pkg/utils/httputils"
	"github.com/codfrm/dns-kit/internal/model/entity/cert_hosting_entity"
	"github.com/codfrm/dns-kit/internal/pkg/code"
	"github.com/codfrm/dns-kit/internal/repository/cdn_repo"
	"github.com/codfrm/dns-kit/internal/repository/cert_hosting_repo"
	"github.com/codfrm/dns-kit/internal/task/queue"
	"github.com/codfrm/dns-kit/internal/task/queue/message"
	"go.uber.org/zap"

	api "github.com/codfrm/dns-kit/internal/api/cert"
)

type HostingSvc interface {
	// HostingList 托管列表
	HostingList(ctx context.Context, req *api.HostingListRequest) (*api.HostingListResponse, error)
	// HostingAdd 添加托管
	HostingAdd(ctx context.Context, req *api.HostingAddRequest) (*api.HostingAddResponse, error)
	// HostingDelete 删除托管
	HostingDelete(ctx context.Context, req *api.HostingDeleteRequest) (*api.HostingDeleteResponse, error)
	// ReDeploy 重新部署
	ReDeploy(ctx context.Context, id int64) error
	// HostingQuery 查询托管
	HostingQuery(ctx context.Context, req *api.HostingQueryRequest) (*api.HostingQueryResponse, error)
}

type hostingSvc struct {
	sync.Mutex
}

var defaultHosting = &hostingSvc{}

func Hosting() HostingSvc {
	return defaultHosting
}

// HostingList 托管列表
func (h *hostingSvc) HostingList(ctx context.Context, req *api.HostingListRequest) (*api.HostingListResponse, error) {
	list, total, err := cert_hosting_repo.CertHosting().FindPage(ctx, req.PageRequest)
	if err != nil {
		return nil, err
	}
	ret := &api.HostingListResponse{
		PageResponse: httputils.PageResponse[*api.HostingItem]{
			List:  make([]*api.HostingItem, 0),
			Total: total,
		},
	}
	for _, v := range list {
		item := &api.HostingItem{
			ID:         v.ID,
			CdnID:      v.CdnID,
			CertID:     v.CertID,
			Status:     v.Status,
			Createtime: v.Createtime,
		}
		switch v.Type {
		case cert_hosting_entity.CertHostingTypeCDN:
			cdn, err := cdn_repo.Cdn().Find(ctx, v.CdnID)
			if err != nil {
				return nil, err
			}
			item.CDN = cdn.Domain
		case cert_hosting_entity.CertHostingTypeProvider:
			deploy, err := deploy.Factory(ctx, v)
			if err != nil {
				return nil, err
			}
			domains, err := deploy.Domains(ctx, v)
			if err != nil {
				return nil, err
			}
			item.CDN = domains[0]
		}
		ret.List = append(ret.List, item)
	}
	return ret, nil
}

// HostingAdd 添加托管
func (h *hostingSvc) HostingAdd(ctx context.Context, req *api.HostingAddRequest) (*api.HostingAddResponse, error) {
	h.Lock()
	defer h.Unlock()
	// 添加入托管
	hosting := &cert_hosting_entity.CertHosting{
		Email:      req.Email,
		Type:       req.Type,
		CdnID:      req.CdnID,
		ProviderID: req.ProviderID,
		Status:     cert_hosting_entity.CertHostingStatusDeploy,
		Createtime: time.Now().Unix(),
		Updatetime: time.Now().Unix(),
	}
	if req.Type == cert_hosting_entity.CertHostingTypeCDN {
		cdn, err := cdn_repo.Cdn().Find(ctx, req.CdnID)
		if err != nil {
			return nil, err
		}
		if err := cdn.Check(ctx); err != nil {
			return nil, err
		}
		// 判断是否已托管
		exists, err := cert_hosting_repo.CertHosting().FindByCDN(ctx, cdn.ID)
		if err != nil {
			return nil, err
		}
		if len(exists) > 0 {
			return nil, i18n.NewError(ctx, code.CertHostingExist)
		}
	} else if req.Type == cert_hosting_entity.CertHostingTypeProvider {
		// 判断厂商是否存在
		provider, err := provider_repo.Provider().Find(ctx, req.ProviderID)
		if err != nil {
			return nil, err
		}
		if err := provider.Check(ctx); err != nil {
			return nil, err
		}
		if err := hosting.SetConfigMap(req.Config); err != nil {
			return nil, err
		}
	} else {
		return nil, i18n.NewError(ctx, code.CertHostingTypeError)
	}
	if err := cert_hosting_repo.CertHosting().Create(ctx, hosting); err != nil {
		return nil, err
	}
	_ = audit.Ctx(ctx).Record("create", zap.Int64("id", hosting.ID),
		zap.Int64("cdn_id", hosting.CdnID), zap.String("email", hosting.Email))
	// 申请部署
	err := Hosting().ReDeploy(ctx, hosting.ID)
	if err != nil {
		logger.Ctx(ctx).Warn("redeploy failed", zap.Error(err))
	}
	return nil, nil
}

// HostingDelete 删除托管
func (h *hostingSvc) HostingDelete(ctx context.Context, req *api.HostingDeleteRequest) (*api.HostingDeleteResponse, error) {
	hosting, err := cert_hosting_repo.CertHosting().Find(ctx, req.ID)
	if err != nil {
		return nil, err
	}
	if err := hosting.Check(ctx); err != nil {
		return nil, err
	}
	if err := cert_hosting_repo.CertHosting().Delete(ctx, hosting.ID); err != nil {
		return nil, err
	}
	return nil, nil
}

func (h *hostingSvc) GetDomains(ctx context.Context, hosting *cert_hosting_entity.CertHosting) ([]string, error) {
	var domains []string
	if hosting.Type == cert_hosting_entity.CertHostingTypeCDN {
		cdn, err := cdn_repo.Cdn().Find(ctx, hosting.CdnID)
		if err != nil {
			return nil, err
		}
		if err := cdn.Check(ctx); err != nil {
			_ = cert_hosting_repo.CertHosting().UpdateStatus(ctx,
				hosting.ID, cert_hosting_entity.CertHostingStatusFail)
			return nil, err
		}
		domains = []string{cdn.Domain}
	} else if hosting.Type == cert_hosting_entity.CertHostingTypeProvider {
		deploy, err := deploy.Factory(ctx, hosting)
		if err != nil {
			return nil, err
		}
		domains, err = deploy.Domains(ctx, hosting)
		if err != nil {
			return nil, err
		}
	}
	return domains, nil
}

// ReDeploy 重新部署
func (h *hostingSvc) ReDeploy(ctx context.Context, id int64) error {
	hosting, err := cert_hosting_repo.CertHosting().Find(ctx, id)
	if err != nil {
		return err
	}
	if err := hosting.Check(ctx); err != nil {
		return err
	}
	if hosting.Status != cert_hosting_entity.CertHostingStatusDeploy {
		return i18n.NewError(ctx, code.CertHostingDeploy)
	}
	domains, err := h.GetDomains(ctx, hosting)
	if err != nil {
		return err
	}
	// 检查域名
	if err := Cert().CheckDomains(ctx, domains); err != nil {
		return err
	}
	defer func() {
		if err != nil {
			_ = cert_hosting_repo.CertHosting().UpdateStatus(ctx,
				hosting.ID, cert_hosting_entity.CertHostingStatusDeployFail)
		}
	}()
	// 发起托管证书更新请求
	cert, err := Cert().Create(ctx, (&api.CreateRequest{
		Email:   hosting.Email,
		Domains: domains,
	}).SetIgnoreMsg(true))
	if err != nil {
		return err
	}
	hosting.CertID = cert.ID
	hosting.Updatetime = time.Now().Unix()
	if err = cert_hosting_repo.CertHosting().Update(ctx, hosting); err != nil {
		return err
	}
	// 这里再发证书创建消息，避免数据还未更新就处理完了
	if err = queue.PublishCertCreate(ctx, &message.CreateCertMessage{ID: cert.ID}); err != nil {
		return err
	}
	return nil
}

// HostingQuery 查询托管
func (h *hostingSvc) HostingQuery(ctx context.Context, req *api.HostingQueryRequest) (*api.HostingQueryResponse, error) {
	cdn, _, err := cdn_repo.Cdn().FindPage(ctx, httputils.PageRequest{Size: 100})
	if err != nil {
		return nil, err
	}
	hosting, _, err := cert_hosting_repo.CertHosting().FindPage(ctx, httputils.PageRequest{Size: 100})
	if err != nil {
		return nil, err
	}
	hostingMap := make(map[int64]bool)
	for _, v := range hosting {
		hostingMap[v.CdnID] = true
	}
	list := make([]*api.HostingQueryItem, 0)
	for _, v := range cdn {
		list = append(list, &api.HostingQueryItem{
			ID:        v.ID,
			Domain:    v.Domain,
			IsManaged: hostingMap[v.ID],
		})
	}
	return &api.HostingQueryResponse{
		List: list,
	}, nil
}
