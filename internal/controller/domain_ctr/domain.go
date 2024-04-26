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
