package dns_svc

import (
	"context"
	"encoding/json"
	"github.com/codfrm/cago/pkg/consts"
	"github.com/codfrm/cago/pkg/i18n"
	"github.com/codfrm/dns-kit/internal/model/entity/dns_provider_entity"
	"github.com/codfrm/dns-kit/internal/pkg/code"
	"github.com/codfrm/dns-kit/internal/repository/dns_provider_repo"
	"time"

	api "github.com/codfrm/dns-kit/internal/api/dns"
)

type ProviderSvc interface {
	// ListProvider 获取dns提供商列表
	ListProvider(ctx context.Context, req *api.ListProviderRequest) (*api.ListProviderResponse, error)
	// CreateProvider 创建dns提供商
	CreateProvider(ctx context.Context, req *api.CreateProviderRequest) (*api.CreateProviderResponse, error)
}

type providerSvc struct {
}

var defaultProvider = &providerSvc{}

func Provider() ProviderSvc {
	return defaultProvider
}

// ListProvider 获取dns提供商列表
func (p *providerSvc) ListProvider(ctx context.Context, req *api.ListProviderRequest) (*api.ListProviderResponse, error) {
	return nil, nil
}

// CreateProvider 创建dns提供商
func (p *providerSvc) CreateProvider(ctx context.Context, req *api.CreateProviderRequest) (*api.CreateProviderResponse, error) {
	manager, err := NewDnsProvider(ctx, req.Platform, req.Secret)
	if err != nil {
		return nil, err
	}
	user, err := manager.UserDetails(ctx)
	if err != nil {
		return nil, i18n.NewError(ctx, code.DNSProviderSecretError)
	}
	// 判断重复添加
	if _, err := dns_provider_repo.DnsProvider().FindByProviderUserId(ctx, user.ID); err == nil {
		return nil, i18n.NewError(ctx, code.DNSProviderExist)
	}
	secretData, err := json.Marshal(req.Secret)
	if err != nil {
		return nil, err
	}
	dnsProvider := &dns_provider_entity.DnsProvider{
		Name:       req.Name,
		UserID:     user.ID,
		Secret:     string(secretData),
		Platform:   req.Platform,
		Status:     consts.ACTIVE,
		Createtime: time.Now().Unix(),
		Updatetime: time.Now().Unix(),
	}
	if err := dns_provider_repo.DnsProvider().Create(ctx, dnsProvider); err != nil {
		return nil, err
	}
	return &api.CreateProviderResponse{}, nil
}
