package domain_ctr

import (
	"context"

	api "github.com/codfrm/dns-kit/internal/api/domain"
	"github.com/codfrm/dns-kit/internal/service/domain_svc"
)

type Domain struct {
}

func NewDomain() *Domain {
	return &Domain{}
}

// List 获取域名列表
func (d *Domain) List(ctx context.Context, req *api.ListRequest) (*api.ListResponse, error) {
	return domain_svc.Domain().List(ctx, req)
}

// Query 查询域名列表
func (d *Domain) Query(ctx context.Context, req *api.QueryRequest) (*api.QueryResponse, error) {
	return domain_svc.Domain().Query(ctx, req)
}

// Add 纳管域名
func (d *Domain) Add(ctx context.Context, req *api.AddRequest) (*api.AddResponse, error) {
	return domain_svc.Domain().Add(ctx, req)
}

// Delete 删除域名
func (d *Domain) Delete(ctx context.Context, req *api.DeleteRequest) (*api.DeleteResponse, error) {
	return domain_svc.Domain().Delete(ctx, req)
}
