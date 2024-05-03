package cloudflare

import (
	"context"
	"strings"

	"github.com/cloudflare/cloudflare-go"
	"github.com/codfrm/dns-kit/pkg/dns"
)

type Manager struct {
	api *cloudflare.API
	rc  *cloudflare.ResourceContainer
}

func NewDNSManager(api *cloudflare.API, rc *cloudflare.ResourceContainer) (dns.Manager, error) {
	return &Manager{
		api: api,
		rc:  rc,
	}, nil
}

func (d *Manager) allDNSRecords(ctx context.Context) ([]cloudflare.DNSRecord, error) {
	resultInfo := cloudflare.ResultInfo{
		Page: 1,
	}
	ret := make([]cloudflare.DNSRecord, 0)
	for {
		records, _, err := d.api.ListDNSRecords(ctx, d.rc, cloudflare.ListDNSRecordsParams{
			ResultInfo: resultInfo,
		})
		if err != nil {
			return nil, err
		}
		ret = append(ret, records...)
		if len(records) == 0 {
			break
		}
		resultInfo.Page++
	}
	return ret, nil
}

func (d *Manager) GetRecordList(ctx context.Context) ([]*dns.Record, error) {
	records, err := d.allDNSRecords(ctx)
	if err != nil {
		return nil, err
	}
	ret := make([]*dns.Record, 0, len(records))
	for _, record := range records {
		ret = append(ret, d.toDNSRecord(record))
	}
	return ret, nil
}

func (d *Manager) toDNSRecord(record cloudflare.DNSRecord) *dns.Record {
	ret := &dns.Record{
		ID:    record.ID,
		Name:  strings.TrimSuffix(record.Name, "."+record.ZoneName),
		Type:  dns.RecordType(record.Type),
		Value: record.Content,
		TTL:   record.TTL,
		Extra: map[string]any{},
	}
	if record.Proxied == nil {
		ret.Extra["proxied"] = false
	} else {
		ret.Extra["proxied"] = *record.Proxied
	}
	return ret
}

func (d *Manager) AddRecord(ctx context.Context, record *dns.Record) error {
	param := cloudflare.CreateDNSRecordParams{
		Type:    string(record.Type),
		Name:    record.Name,
		Content: record.Value,
		TTL:     record.TTL,
	}
	param.Proxied = new(bool)
	if d.isProxied(record) {
		*param.Proxied = true
	} else {
		*param.Proxied = false
	}
	if record.TTL == 0 {
		record.TTL = 1
	}
	resp, err := d.api.CreateDNSRecord(ctx, d.rc, param)
	if err != nil {
		return err
	}
	record.ID = resp.ID
	return nil
}

func (d *Manager) UpdateRecord(ctx context.Context, recordId string, record *dns.Record) error {
	param := cloudflare.UpdateDNSRecordParams{
		Type:    string(record.Type),
		Name:    record.Name,
		Content: record.Value,
		ID:      recordId,
		TTL:     record.TTL,
	}
	param.Proxied = new(bool)
	if d.isProxied(record) {
		*param.Proxied = true
	} else {
		*param.Proxied = false
	}
	_, err := d.api.UpdateDNSRecord(ctx, d.rc, param)
	if err != nil {
		return err
	}
	return nil
}

func (d *Manager) isProxied(record *dns.Record) bool {
	ret, ok := record.Extra["proxied"].(bool)
	if !ok {
		return false
	}
	return ret
}

func (d *Manager) DelRecord(ctx context.Context, recordId string) error {
	return d.api.DeleteDNSRecord(ctx, d.rc, recordId)
}

func (d *Manager) ExtraFields() []*dns.Extra {
	return []*dns.Extra{{
		Key:       "proxied",
		Title:     "代理",
		FieldType: dns.FieldTypeSwitch,
		Default:   true,
	}}
}
