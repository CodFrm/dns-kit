package provider_ctr

import (
	"context"

	api "github.com/codfrm/dns-kit/internal/api/provider"
	"github.com/codfrm/dns-kit/internal/service/provider_svc"
)

type Provider struct {
}

func NewProvider() *Provider {
	return &Provider{}
}

// ListProvider 获取供应商列表
func (p *Provider) ListProvider(ctx context.Context, req *api.ListProviderRequest) (*api.ListProviderResponse, error) {
	return provider_svc.Provider().ListProvider(ctx, req)
}

// CreateProvider 创建供应商
func (p *Provider) CreateProvider(ctx context.Context, req *api.CreateProviderRequest) (*api.CreateProviderResponse, error) {
	return provider_svc.Provider().CreateProvider(ctx, req)
}

// UpdateProvider 更新供应商
func (p *Provider) UpdateProvider(ctx context.Context, req *api.UpdateProviderRequest) (*api.UpdateProviderResponse, error) {
	return provider_svc.Provider().UpdateProvider(ctx, req)
}

// DeleteProvider 删除供应商
func (p *Provider) DeleteProvider(ctx context.Context, req *api.DeleteProviderRequest) (*api.DeleteProviderResponse, error) {
	return provider_svc.Provider().DeleteProvider(ctx, req)
}
