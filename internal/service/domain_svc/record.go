package domain_svc

import (
	"context"
	"sync"

	"github.com/codfrm/cago/pkg/i18n"
	"github.com/codfrm/cago/pkg/iam/audit"
	"github.com/codfrm/dns-kit/internal/pkg/code"
	"github.com/codfrm/dns-kit/internal/repository/domain_repo"
	"github.com/codfrm/dns-kit/pkg/platform"
	"go.uber.org/zap"

	api "github.com/codfrm/dns-kit/internal/api/domain"
)

type RecordSvc interface {
	// RecordList 获取记录列表
	RecordList(ctx context.Context, req *api.RecordListRequest) (*api.RecordListResponse, error)
	// CreateRecord 创建记录
	CreateRecord(ctx context.Context, req *api.CreateRecordRequest) (*api.CreateRecordResponse, error)
	// UpdateRecord 更新记录
	UpdateRecord(ctx context.Context, req *api.UpdateRecordRequest) (*api.UpdateRecordResponse, error)
	// DeleteRecord 删除记录
	DeleteRecord(ctx context.Context, req *api.DeleteRecordRequest) (*api.DeleteRecordResponse, error)
}

type recordSvc struct {
	sync.Mutex
}

var defaultRecord = &recordSvc{}

func Record() RecordSvc {
	return defaultRecord
}

// RecordList 获取记录列表
func (r *recordSvc) RecordList(ctx context.Context, req *api.RecordListRequest) (*api.RecordListResponse, error) {
	domain, err := domain_repo.Domain().Find(ctx, req.ID)
	if err != nil {
		return nil, err
	}
	dnsManager, err := domain.DnsManager(ctx)
	if err != nil {
		return nil, err
	}
	list, err := dnsManager.GetRecordList(ctx)
	if err != nil {
		return nil, err
	}
	return &api.RecordListResponse{
		List:  list,
		Extra: dnsManager.ExtraFields(),
	}, nil
}

// CreateRecord 创建记录
func (r *recordSvc) CreateRecord(ctx context.Context, req *api.CreateRecordRequest) (*api.CreateRecordResponse, error) {
	domain, err := domain_repo.Domain().Find(ctx, req.ID)
	if err != nil {
		return nil, err
	}
	dnsManager, err := domain.DnsManager(ctx)
	if err != nil {
		return nil, err
	}
	if err := dnsManager.AddRecord(ctx, &platform.Record{
		Type:  req.Type,
		Name:  req.Name,
		Value: req.Value,
		TTL:   req.TTL,
		Extra: req.Extra,
	}); err != nil {
		return nil, i18n.NewError(ctx, code.RecordCreateFailed, err)
	}
	_ = audit.Ctx(ctx).Record("create", zap.String("record_name", req.Name))
	return nil, nil
}

// UpdateRecord 更新记录
func (r *recordSvc) UpdateRecord(ctx context.Context, req *api.UpdateRecordRequest) (*api.UpdateRecordResponse, error) {
	domain, err := domain_repo.Domain().Find(ctx, req.ID)
	if err != nil {
		return nil, err
	}
	dnsManager, err := domain.DnsManager(ctx)
	if err != nil {
		return nil, err
	}
	if err := dnsManager.UpdateRecord(ctx, req.RecordID, &platform.Record{
		Type:  req.Type,
		Name:  req.Name,
		Value: req.Value,
		TTL:   req.TTL,
		Extra: req.Extra,
	}); err != nil {
		return nil, i18n.NewError(ctx, code.RecordUpdateFailed, err)
	}
	_ = audit.Ctx(ctx).Record("update", zap.String("record_id", req.RecordID))
	return nil, nil
}

// DeleteRecord 删除记录
func (r *recordSvc) DeleteRecord(ctx context.Context, req *api.DeleteRecordRequest) (*api.DeleteRecordResponse, error) {
	domain, err := domain_repo.Domain().Find(ctx, req.ID)
	if err != nil {
		return nil, err
	}
	dnsManager, err := domain.DnsManager(ctx)
	if err != nil {
		return nil, err
	}
	err = dnsManager.DelRecord(ctx, req.RecordID)
	if err != nil {
		return nil, err
	}
	_ = audit.Ctx(ctx).Record("delete", zap.String("record_id", req.RecordID))
	return nil, nil
}
