package dnspod

import (
	"context"
	"fmt"
	"strconv"

	"github.com/codfrm/dns-kit/pkg/dns"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	dnspod "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dnspod/v20210323"
)

type Manager struct {
	api *dnspod.Client
	rc  *dnspod.DomainInfo
}

func (m Manager) GetRecordList(ctx context.Context) ([]*dns.Record, error) {
	// 实例化一个请求对象,每个接口都会对应一个request对象
	request := dnspod.NewDescribeRecordListRequest()
	request.SetContext(ctx)
	request.Domain = common.StringPtr(*m.rc.Domain)

	// 返回的resp是一个DescribeRecordListResponse的实例，与请求对象对应
	response, err := m.api.DescribeRecordList(request)
	if _, ok := err.(*errors.TencentCloudSDKError); ok {
		fmt.Printf("An API error has returned: %s", err)
		return nil, err
	}
	if err != nil {
		return nil, err
	}
	records := make([]*dns.Record, 0, len(response.Response.RecordList))
	for _, record := range response.Response.RecordList {
		records = append(records, &dns.Record{
			ID:    strconv.FormatUint(*record.RecordId, 10),
			Type:  dns.RecordType(*record.Type),
			Name:  *record.Name,
			Value: *record.Value,
			TTL:   int(*record.TTL),
			Extra: map[string]any{"line": *record.Line},
		})
	}
	return records, nil
}

func (m Manager) AddRecord(ctx context.Context, record *dns.Record) error {
	// 实例化一个请求对象,每个接口都会对应一个request对象
	request := dnspod.NewCreateRecordRequest()
	atoi, err := strconv.ParseUint(record.ID, 10, 64)
	if err != nil {
		return err
	}
	//域名id
	request.DomainId = common.Uint64Ptr(atoi)
	//记录类型，通过 API 记录类型获得，大写英文，比如：A 。
	request.RecordType = common.StringPtr(string(record.Type))
	//主机记录，如 www，如果不传，默认为 @。
	//示例值：www
	request.SubDomain = common.StringPtr(record.Name)
	//记录值，如 IP : 200.200.200.200， CNAME : cname.dnspod.com.， MX : mail.dnspod.com.。
	//示例值：200.200.200.200
	request.Value = common.StringPtr(record.Value)
	RecordLine, ok := record.Extra["RecordLine"].(string)
	if !ok {
		return nil
	}
	request.RecordLine = common.StringPtr(RecordLine)

	// 返回的resp是一个CreateRecordResponse的实例，与请求对象对应
	_, err = m.api.CreateRecord(request)
	if _, ok := err.(*errors.TencentCloudSDKError); ok {
		return err

	}
	if err != nil {
		return err
	}
	return nil
}

func (m Manager) UpdateRecord(ctx context.Context, recordId string, record *dns.Record) error {
	// 实例化一个请求对象,每个接口都会对应一个request对象
	request := dnspod.NewModifyRecordRequest()
	request.SetContext(ctx)
	request.Domain = common.StringPtr(*m.rc.Domain)
	request.RecordType = common.StringPtr(string(record.Type))
	//额外字段线路
	RecordLine := record.Extra["RecordLine"].(string)

	request.RecordLine = common.StringPtr(RecordLine)
	request.Value = common.StringPtr(record.Value)
	//类型转换
	parseUint, err := strconv.ParseUint(recordId, 10, 64)
	if err != nil {
		return err
	}
	request.RecordId = common.Uint64Ptr(parseUint)
	request.SubDomain = common.StringPtr(record.Name)

	// 返回的resp是一个ModifyRecordResponse的实例，与请求对象对应
	_, err = m.api.ModifyRecord(request)
	if _, ok := err.(*errors.TencentCloudSDKError); ok {
		fmt.Printf("An API error has returned: %s", err)
		return err
	}
	if err != nil {
		return err
	}
	return nil
}

func (m Manager) DelRecord(ctx context.Context, recordId string) error {
	// 实例化一个请求对象,每个接口都会对应一个request对象
	request := dnspod.NewDeleteRecordRequest()

	request.Domain = common.StringPtr(*m.rc.Domain)
	parseUint, err := strconv.ParseUint(recordId, 10, 64)
	if err != nil {
		return err
	}
	request.RecordId = common.Uint64Ptr(parseUint)

	// 返回的resp是一个DeleteRecordResponse的实例，与请求对象对应
	_, err = m.api.DeleteRecord(request)
	if _, ok := err.(*errors.TencentCloudSDKError); ok {
		fmt.Printf("An API error has returned: %s", err)
		return err
	}
	if err != nil {
		return err
	}
	return nil
}

func (m Manager) ExtraFields() []*dns.Extra {
	//额外字段:线路
	return []*dns.Extra{{
		Key:       "line",
		Title:     "线路",
		FieldType: dns.FieldTypeSelect,
		Options: []string{
			"默认",
			"电信", "联通", "移动", "铁通", "广电", "教育网", "境内", "境外",
			"百度", "谷歌", "有道", "必应", "搜狗", "奇虎", "搜索引擎",
		},
		Default: "默认",
	}}
}

func NewDNSManager(api *dnspod.Client, rc *dnspod.DomainInfo) (dns.Manager, error) {
	return &Manager{
		api: api,
		rc:  rc,
	}, nil
}
