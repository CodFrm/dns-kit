package domain_svc

import (
	"context"
	"github.com/codfrm/dns-kit/internal/repository/domain_repo"

	api "github.com/codfrm/dns-kit/internal/api/domain"
)

type DomainSvc interface {
	// List 获取域名列表
	List(ctx context.Context, req *api.ListRequest) (*api.ListResponse, error)
}

type domainSvc struct {
}

var defaultDomain = &domainSvc{}

func Domain() DomainSvc {
	return defaultDomain
}

// List 获取域名列表
func (d *domainSvc) List(ctx context.Context, req *api.ListRequest) (*api.ListResponse, error) {
	domain_repo.Domain().FindPage()
	return nil, nil
}
