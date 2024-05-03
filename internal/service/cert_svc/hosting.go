package cert_svc

import (
	"context"
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
	"sync"
	"time"

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
		cdn, err := cdn_repo.Cdn().Find(ctx, v.CdnID)
		if err != nil {
			return nil, err
		}
		ret.List = append(ret.List, &api.HostingItem{
			ID:         v.ID,
			CdnID:      v.CdnID,
			CDN:        cdn.Domain,
			CertID:     v.CertID,
			Status:     v.Status,
			Createtime: v.Createtime,
		})
	}
	return ret, nil
}

// HostingAdd 添加托管
func (h *hostingSvc) HostingAdd(ctx context.Context, req *api.HostingAddRequest) (*api.HostingAddResponse, error) {
	h.Lock()
	defer h.Unlock()
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
	// 添加入托管
	hosting := &cert_hosting_entity.CertHosting{
		Email:      req.Email,
		CdnID:      req.CdnID,
		Status:     cert_hosting_entity.CertHostingStatusActive,
		Createtime: time.Now().Unix(),
		Updatetime: time.Now().Unix(),
	}
	if err := cert_hosting_repo.CertHosting().Create(ctx, hosting); err != nil {
		return nil, err
	}
	_ = audit.Ctx(ctx).Record("create", zap.Int64("id", hosting.ID),
		zap.Int64("cdn_id", hosting.CdnID), zap.String("email", hosting.Email))
	// 申请部署
	err = Hosting().ReDeploy(ctx, hosting.ID)
	if err != nil {
		logger.Ctx(ctx).Warn("redeploy failed", zap.Error(err))
	}
	return nil, nil
}

// HostingDelete 删除托管
func (h *hostingSvc) HostingDelete(ctx context.Context, req *api.HostingDeleteRequest) (*api.HostingDeleteResponse, error) {
	return nil, nil
}

// ReDeploy 重新部署
func (h *hostingSvc) ReDeploy(ctx context.Context, id int64) error {
	h.Lock()
	defer h.Unlock()
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
	cdn, err := cdn_repo.Cdn().Find(ctx, hosting.CdnID)
	if err != nil {
		return err
	}
	if err := cdn.Check(ctx); err != nil {
		_ = cert_hosting_repo.CertHosting().UpdateStatus(ctx,
			hosting.ID, cert_hosting_entity.CertHostingStatusFail)
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
		Domains: []string{cdn.Domain},
	}).SetIgnoreMsg(true))
	if err != nil {
		return err
	}
	hosting.Status = cert_hosting_entity.CertHostingStatusDeploy
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
