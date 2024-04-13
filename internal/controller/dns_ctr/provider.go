package dns_ctr

import (
	"context"

	api "github.com/codfrm/dns-kit/internal/api/dns"
	"github.com/codfrm/dns-kit/internal/service/dns_svc"
)

type Provider struct {
}

func NewProvider() *Provider {
	return &Provider{}
}

// ListProvider 获取dns提供商列表
func (p *Provider) ListProvider(ctx context.Context, req *api.ListProviderRequest) (*api.ListProviderResponse, error) {
	return dns_svc.Provider().ListProvider(ctx, req)
}

// CreateProvider 创建dns提供商
func (p *Provider) CreateProvider(ctx context.Context, req *api.CreateProviderRequest) (*api.CreateProviderResponse, error) {
	return dns_svc.Provider().CreateProvider(ctx, req)
}
