package cdn_ctr

import (
	"context"

	api "github.com/codfrm/dns-kit/internal/api/cdn"
	"github.com/codfrm/dns-kit/internal/service/cdn_svc"
)

type Cdn struct {
}

func NewCdn() *Cdn {
	return &Cdn{}
}

// Query 查询cdn
func (c *Cdn) Query(ctx context.Context, req *api.QueryRequest) (*api.QueryResponse, error) {
	return cdn_svc.Cdn().Query(ctx, req)
}

// List 获取纳管的cdn列表
func (c *Cdn) List(ctx context.Context, req *api.ListRequest) (*api.ListResponse, error) {
	return cdn_svc.Cdn().List(ctx, req)
}

// Add 添加cdn进入纳管
func (c *Cdn) Add(ctx context.Context, req *api.AddRequest) (*api.AddResponse, error) {
	return cdn_svc.Cdn().Add(ctx, req)
}

// Delete 删除cdn
func (c *Cdn) Delete(ctx context.Context, req *api.DeleteRequest) (*api.DeleteResponse, error) {
	return cdn_svc.Cdn().Delete(ctx, req)
}
