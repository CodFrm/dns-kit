package provider_svc

import (
	"context"

	"github.com/codfrm/cago/pkg/i18n"
	"github.com/codfrm/dns-kit/internal/model/entity/provider_entity"
	"github.com/codfrm/dns-kit/internal/pkg/code"
	"github.com/codfrm/dns-kit/pkg/dns"
	"github.com/codfrm/dns-kit/pkg/dns/provider/cloudflare"
	"github.com/codfrm/dns-kit/pkg/dns/provider/dnspod"
)

func NewProvider(ctx context.Context, platform provider_entity.Platform, secret map[string]string) (manager dns.DomainManager, err error) {
	switch platform {
	case provider_entity.PlatformCloudflare:
		manager, err = cloudflare.NewCloudflare(secret["key"], secret["email"])
	case provider_entity.PlatformDnsPod:
		manager, err = dnspod.NewDnsPod(secret["secret_id"], secret["secret_key"])
	default:
		return nil, i18n.NewError(ctx, code.DNSProviderNotSupport)
	}
	if err != nil {
		return nil, err
	}
	return manager, nil
}
