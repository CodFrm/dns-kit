package cert_ctr

import (
	"context"

	api "github.com/codfrm/dns-kit/internal/api/cert"
	"github.com/codfrm/dns-kit/internal/service/cert_svc"
)

type Hosting struct {
}

func NewHosting() *Hosting {
	return &Hosting{}
}

// HostingList 托管列表
func (h *Hosting) HostingList(ctx context.Context, req *api.HostingListRequest) (*api.HostingListResponse, error) {
	return cert_svc.Hosting().HostingList(ctx, req)
}

// HostingAdd 添加托管
func (h *Hosting) HostingAdd(ctx context.Context, req *api.HostingAddRequest) (*api.HostingAddResponse, error) {
	return cert_svc.Hosting().HostingAdd(ctx, req)
}

// HostingDelete 删除托管
func (h *Hosting) HostingDelete(ctx context.Context, req *api.HostingDeleteRequest) (*api.HostingDeleteResponse, error) {
	return cert_svc.Hosting().HostingDelete(ctx, req)
}

// HostingQuery 查询托管
func (h *Hosting) HostingQuery(ctx context.Context, req *api.HostingQueryRequest) (*api.HostingQueryResponse, error) {
	return cert_svc.Hosting().HostingQuery(ctx, req)
}
