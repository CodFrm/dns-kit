package tencent

import (
	"context"
	"github.com/codfrm/dns-kit/pkg/platform"
	cdn "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdn/v20180606"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
)

// GetCDNList 获取cdn列表
func (t *Tencent) GetCDNList(ctx context.Context) ([]*platform.CDNItem, error) {
	req := cdn.NewDescribeDomainsRequest()
	req.SetContext(ctx)
	resp, err := t.cdnApi.DescribeDomains(req)
	if err != nil {
		return nil, err
	}
	list := make([]*platform.CDNItem, 0)
	for _, v := range resp.Response.Domains {
		list = append(list, &platform.CDNItem{
			ID:     *v.ResourceId,
			Domain: *v.Domain,
		})
	}
	return list, nil
}

func (t *Tencent) GetCDNDetail(ctx context.Context, domain *platform.CDNItem) (*platform.CDNItem, error) {
	req := cdn.NewDescribeDomainsRequest()
	req.SetContext(ctx)
	req.Filters = []*cdn.DomainFilter{{
		Name:  common.StringPtr("resourceId"),
		Value: common.StringPtrs([]string{domain.ID}),
		Fuzzy: common.BoolPtr(false),
	}}
	resp, err := t.cdnApi.DescribeDomains(req)
	if err != nil {
		return nil, err
	}
	if len(resp.Response.Domains) == 0 {
		return nil, nil
	}
	return domain, nil
}

// SetCDNHttpsCert 设置cdn https证书
func (t *Tencent) SetCDNHttpsCert(ctx context.Context, domain *platform.CDNItem, cert, key string) error {
	req := cdn.NewUpdateDomainConfigRequest()
	req.SetContext(ctx)
	req.Domain = common.StringPtr(domain.Domain)
	//req.ProjectId = common.Int64Ptr(t.projectId)
	req.Https = &cdn.Https{
		CertInfo: &cdn.ServerCert{
			Certificate: common.StringPtr(cert),
			PrivateKey:  common.StringPtr(key),
		},
	}
	_, err := t.cdnApi.UpdateDomainConfig(req)
	if err != nil {
		return err
	}
	return nil
}
