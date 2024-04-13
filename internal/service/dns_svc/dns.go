package dns_svc

import (
	"context"

	api "github.com/codfrm/dns-kit/internal/api/dns"
)

type DnsSvc interface {
	// List 获取dns列表
	List(ctx context.Context, req *api.ListRequest) (*api.ListResponse, error)
}

type dnsSvc struct {
}

var defaultDns = &dnsSvc{}

func Dns() DnsSvc {
	return defaultDns
}

// List 获取dns列表
func (d *dnsSvc) List(ctx context.Context, req *api.ListRequest) (*api.ListResponse, error) {
	return nil, nil
}
