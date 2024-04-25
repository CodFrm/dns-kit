package dnspod

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/codfrm/cago/pkg/logger"
	"github.com/codfrm/dns-kit/pkg/dns"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	errors2 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	dnspod "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dnspod/v20210323"
)

type DnsPod struct {
	api *dnspod.Client
}

func (t *DnsPod) GetDomainList(ctx context.Context) ([]*dns.Domain, error) {
	request := dnspod.NewDescribeDomainListRequest()
	// 返回的resp是一个DescribeDomainListResponse的实例，与请求对象对应
	response, err := t.api.DescribeDomainList(request)
	var tencentCloudSDKError *errors2.TencentCloudSDKError
	if errors.As(err, &tencentCloudSDKError) {
		logger.Default().Error(fmt.Sprintf("An API error has returned: %s", err))
		return nil, err
	}
	if err != nil {
		logger.Default().Error("请求DnsPod失败")
		return nil, err
	}
	result := make([]*dns.Domain, 0)
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

func (t *DnsPod) UserDetails(ctx context.Context) (*dns.User, error) {
	resp, err := t.api.DescribeUserDetailWithContext(ctx, dnspod.NewDescribeUserDetailRequest())
	if err != nil {
		return nil, err
	}
	return &dns.User{
		ID:       strconv.FormatInt(*resp.Response.UserInfo.Id, 10),
		Username: *resp.Response.UserInfo.Nick,
	}, nil
}

func NewDnsPod(SecretId, SecretKey string) (dns.DomainManager, error) {
	credential := common.NewCredential(
		SecretId,
		SecretKey,
	)
	// 实例化一个client选项，可选的，没有特殊需求可以跳过
	cpf := profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = "dnspod.tencentcloudapi.com"
	// 实例化要请求产品的client对象,clientProfile是可选的
	client, _ := dnspod.NewClient(credential, "", cpf)
	return &DnsPod{api: client}, nil

}
