package dns_ctr

import (
	"context"

	api "github.com/codfrm/dns-kit/internal/api/dns"
	"github.com/codfrm/dns-kit/internal/service/dns_svc"
)

type Dns struct {
}

func NewDns() *Dns {
	return &Dns{}
}

// List 获取dns列表
func (d *Dns) List(ctx context.Context, req *api.ListRequest) (*api.ListResponse, error) {
	return dns_svc.Dns().List(ctx, req)
}
