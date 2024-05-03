package cloudflare

import (
	"context"

	"github.com/cloudflare/cloudflare-go"

	"github.com/codfrm/dns-kit/pkg/platform"
)

type Cloudflare struct {
	api *cloudflare.API
}

func NewCloudflare(token string) (platform.DomainManager, error) {
	api, err := cloudflare.NewWithAPIToken(token)
	if err != nil {
		return nil, err
	}
	return &Cloudflare{
		api: api,
	}, nil
}

func (c *Cloudflare) GetDomainList(ctx context.Context) ([]*platform.Domain, error) {
	zones, err := c.api.ListZonesContext(ctx)
	if err != nil {
		return nil, err
	}
	ret := make([]*platform.Domain, 0, len(zones.Result))
	for _, zone := range zones.Result {
		ret = append(ret, c.toDomain(zone))
	}
	return ret, nil
}

func (c *Cloudflare) toDomain(zone cloudflare.Zone) *platform.Domain {
	return &platform.Domain{
		ID:     zone.ID,
		Domain: zone.Name,
	}
}

func (c *Cloudflare) BuildDNSManager(ctx context.Context, domain *platform.Domain) (platform.DNSManager, error) {
	return NewDNSManager(c.api, cloudflare.ZoneIdentifier(domain.ID))
}

func (c *Cloudflare) UserDetails(ctx context.Context) (*platform.User, error) {
	user, err := c.api.UserDetails(ctx)
	if err != nil {
		return nil, err
	}
	return &platform.User{
		ID:       user.ID,
		Username: user.Username,
	}, nil
}
