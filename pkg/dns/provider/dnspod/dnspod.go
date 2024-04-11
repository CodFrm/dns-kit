package dnspod

import (
	"context"
	"fmt"
	"github.com/codfrm/cago/pkg/logger"
	"github.com/codfrm/dns-kit/pkg/dns"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	dnspod "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dnspod/v20210323"
	"strconv"
)

type DnsPod struct {
	api *dnspod.Client
}

func (t *DnsPod) GetDomainList(ctx context.Context) ([]*dns.Domain, error) {
	request := dnspod.NewDescribeDomainListRequest()
	// 返回的resp是一个DescribeDomainListResponse的实例，与请求对象对应
	response, err := t.api.DescribeDomainList(request)
	if _, ok := err.(*errors.TencentCloudSDKError); ok {
		logger.Default().Error(fmt.Sprintf("An API error has returned: %s", err))
		return nil, err
	}
	if err != nil {
		logger.Default().Error("请求DnsPod失败")
		return nil, err
	}
	var result []*dns.Domain
	for _, item := range response.Response.DomainList {
		result = append(result, &dns.Domain{
			ID:     strconv.FormatUint(*item.DomainId, 10),
			Domain: *item.Name,
		})
	}
	return result, nil
}

// BuildDNSManager 创建域名管理器
func (t *DnsPod) BuildDNSManager(ctx context.Context, domain *dns.Domain) (dns.Manager, error) {
	parseUint, err := strconv.ParseUint(domain.ID, 10, 64)
	if err != nil {
		return nil, err
	}
	return NewDNSManager(t.api, &dnspod.DomainInfo{DomainId: &parseUint, Domain: &domain.Domain})
}

func NewTencentCloud(SecretId, SecretKey string) (dns.DomainManager, error) {
	credential := common.NewCredential(
		SecretKey,
		SecretKey,
	)
	// 实例化一个client选项，可选的，没有特殊需求可以跳过
	cpf := profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = "dnspod.tencentcloudapi.com"
	// 实例化要请求产品的client对象,clientProfile是可选的
	client, _ := dnspod.NewClient(credential, "", cpf)
	return &DnsPod{api: client}, nil

}

// GetDomainLIst 获取域名列表
func (t *DnsPod) GetDomainLIst(ctx context.Context) ([]*dns.Domain, error) {
	request := dnspod.NewDescribeDomainListRequest()
	// 返回的resp是一个DescribeDomainListResponse的实例，与请求对象对应
	response, err := t.api.DescribeDomainList(request)
	if _, ok := err.(*errors.TencentCloudSDKError); ok {
		fmt.Printf("An API error has returned: %s", err)
		return nil, err
	}
	if err != nil {
		return nil, err
	}
	ret := make([]*dns.Domain, 0, len(response.Response.DomainList))
	// 输出json格式的字符串回包
	for _, item := range response.Response.DomainList {
		domain := dns.Domain{
			ID:     strconv.FormatUint(*item.DomainId, 10),
			Domain: *item.Name,
		}
		ret = append(ret, &domain)
	}
	return ret, nil
}
