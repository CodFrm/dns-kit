package dns_svc

import (
	"context"
	"github.com/codfrm/cago/pkg/i18n"
	"github.com/codfrm/dns-kit/internal/model/entity/dns_provider_entity"
	"github.com/codfrm/dns-kit/internal/pkg/code"
	"github.com/codfrm/dns-kit/pkg/dns"
	"github.com/codfrm/dns-kit/pkg/dns/provider/cloudflare"
	"github.com/codfrm/dns-kit/pkg/dns/provider/dnspod"
)

func NewDnsProvider(ctx context.Context, platform dns_provider_entity.DnsPlatform, secret map[string]string) (manager dns.DomainManager, err error) {
	switch platform {
	case dns_provider_entity.DnsPlatformCloudflare:
		manager, err = cloudflare.NewCloudflare(secret["key"], secret["email"])
	case dns_provider_entity.DnsPlatformDnsPod:
		manager, err = dnspod.NewDnsPod(secret["secret_id"], secret["secret_key"])
	default:
		return nil, i18n.NewError(ctx, code.DNSProviderNotSupport)
	}
	if err != nil {
		return nil, err
	}
	return manager, nil
}
