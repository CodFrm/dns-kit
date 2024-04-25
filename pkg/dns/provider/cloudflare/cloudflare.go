package cloudflare

import (
	"context"

	"github.com/cloudflare/cloudflare-go"

	"github.com/codfrm/dns-kit/pkg/dns"
)

type Cloudflare struct {
	api *cloudflare.API
}

func NewCloudflare(key string, email string) (dns.DomainManager, error) {
	api, err := cloudflare.New(key, email)
	if err != nil {
		return nil, err
	}
	return &Cloudflare{
		api: api,
	}, nil
}

func (c *Cloudflare) GetDomainList(ctx context.Context) ([]*dns.Domain, error) {
	zones, err := c.api.ListZonesContext(ctx)
	if err != nil {
		return nil, err
	}
	ret := make([]*dns.Domain, 0, len(zones.Result))
	for _, zone := range zones.Result {
		ret = append(ret, c.toDomain(zone))
	}
	return ret, nil
}

func (c *Cloudflare) toDomain(zone cloudflare.Zone) *dns.Domain {
	return &dns.Domain{
		ID:     zone.ID,
		Domain: zone.Name,
	}
}

func (c *Cloudflare) BuildDNSManager(ctx context.Context, domain *dns.Domain) (dns.Manager, error) {
	return NewDNSManager(c.api, cloudflare.ZoneIdentifier(domain.ID))
}

func (c *Cloudflare) UserDetails(ctx context.Context) (*dns.User, error) {
	user, err := c.api.UserDetails(ctx)
	if err != nil {
		return nil, err
	}
	return &dns.User{
		ID:       user.ID,
		Username: user.Username,
	}, nil
}
