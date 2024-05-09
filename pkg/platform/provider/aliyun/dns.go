package aliyun

import (
	"context"
	alidns20150109 "github.com/alibabacloud-go/alidns-20150109/v4/client"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/codfrm/dns-kit/pkg/platform"
)

type DNSManager struct {
	dnsClient *alidns20150109.Client
	domain    *platform.Domain
}

func NewDNSManager(dnsClient *alidns20150109.Client, domain *platform.Domain) (platform.DNSManager, error) {
	return &DNSManager{
		dnsClient: dnsClient,
		domain:    domain,
	}, nil
}

func (d *DNSManager) GetRecordList(ctx context.Context) ([]*platform.Record, error) {
	resp, err := d.dnsClient.DescribeDomainRecords(&alidns20150109.DescribeDomainRecordsRequest{
		DomainName: tea.String(d.domain.Domain),
		PageSize:   tea.Int64(500),
	})
	if err != nil {
		return nil, err
	}
	ret := make([]*platform.Record, 0, len(resp.Body.DomainRecords.Record))
	for _, v := range resp.Body.DomainRecords.Record {
		ret = append(ret, &platform.Record{
			ID:    *v.RecordId,
			Type:  platform.RecordType(*v.Type),
			Name:  *v.RR,
			Value: *v.Value,
			TTL:   int(*v.TTL),
			Extra: nil,
		})
	}
	return ret, nil
}

func (d *DNSManager) AddRecord(ctx context.Context, record *platform.Record) error {
	req := &alidns20150109.AddDomainRecordRequest{
		DomainName: tea.String(d.domain.Domain),
		RR:         tea.String(record.Name),
		Type:       tea.String(string(record.Type)),
		Value:      tea.String(record.Value),
	}
	if record.TTL == 0 {
		req.TTL = nil
	}
	_, err := d.dnsClient.AddDomainRecord(req)
	if err != nil {
		return err
	}
	return nil
}

func (d *DNSManager) UpdateRecord(ctx context.Context, recordId string, record *platform.Record) error {
	req := &alidns20150109.UpdateDomainRecordRequest{
		RecordId: tea.String(recordId),
		RR:       tea.String(record.Name),
		Type:     tea.String(string(record.Type)),
		Value:    tea.String(record.Value),
	}
	if record.TTL == 0 {
		req.TTL = nil
	}
	_, err := d.dnsClient.UpdateDomainRecord(req)
	if err != nil {
		return err
	}
	return nil
}

func (d *DNSManager) DelRecord(ctx context.Context, recordId string) error {
	_, err := d.dnsClient.DeleteDomainRecord(&alidns20150109.DeleteDomainRecordRequest{
		RecordId: tea.String(recordId),
	})
	if err != nil {
		return err
	}
	return nil
}

func (d *DNSManager) ExtraFields() []*platform.Extra {
	return []*platform.Extra{}
}
