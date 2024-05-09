package aliyun

import (
	"context"
	alidns20150109 "github.com/alibabacloud-go/alidns-20150109/v4/client"
	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	sts20150401 "github.com/alibabacloud-go/sts-20150401/v2/client"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/codfrm/dns-kit/pkg/platform"
)

type Aliyun struct {
	config    *openapi.Config
	dnsClient *alidns20150109.Client
}

func NewAliyun(accessKeyId, accessKeySecret string) (*Aliyun, error) {
	config := &openapi.Config{
		AccessKeyId:     tea.String(accessKeyId),
		AccessKeySecret: tea.String(accessKeySecret),
	}

	// Endpoint 请参考 https://api.aliyun.com/product/Alidns
	config.Endpoint = tea.String("alidns.cn-hangzhou.aliyuncs.com")
	dnsClient, err := alidns20150109.NewClient(config)
	if err != nil {
		return nil, err
	}
	return &Aliyun{dnsClient: dnsClient, config: config}, nil
}

func (a *Aliyun) UserDetails(ctx context.Context) (*platform.User, error) {
	config := &openapi.Config{
		AccessKeyId:     a.config.AccessKeyId,
		AccessKeySecret: a.config.AccessKeySecret,
	}
	config.Endpoint = tea.String("sts.cn-guangzhou.aliyuncs.com")
	result, err := sts20150401.NewClient(config)
	if err != nil {
		return nil, err
	}
	resp, err := result.GetCallerIdentity()
	if err != nil {
		return nil, err
	}

	return &platform.User{
		ID:       *resp.Body.UserId,
		Username: *resp.Body.AccountId,
	}, nil
}

func (a *Aliyun) GetDomainList(ctx context.Context) ([]*platform.Domain, error) {
	resp, err := a.dnsClient.DescribeDomains(&alidns20150109.DescribeDomainsRequest{})
	if err != nil {
		return nil, err
	}
	ret := make([]*platform.Domain, 0, len(resp.Body.Domains.Domain))
	for _, domain := range resp.Body.Domains.Domain {
		ret = append(ret, &platform.Domain{
			ID:     *domain.DomainId,
			Domain: *domain.DomainName,
		})
	}
	return ret, nil
}

func (a *Aliyun) BuildDNSManager(ctx context.Context, domain *platform.Domain) (platform.DNSManager, error) {
	return NewDNSManager(a.dnsClient, domain)
}
