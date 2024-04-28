package domain_ctr

import (
	"context"

	api "github.com/codfrm/dns-kit/internal/api/domain"
	"github.com/codfrm/dns-kit/internal/service/domain_svc"
)

type Record struct {
}

func NewRecord() *Record {
	return &Record{}
}

// RecordList 获取记录列表
func (r *Record) RecordList(ctx context.Context, req *api.RecordListRequest) (*api.RecordListResponse, error) {
	return domain_svc.Record().RecordList(ctx, req)
}

// CreateRecord 创建记录
func (r *Record) CreateRecord(ctx context.Context, req *api.CreateRecordRequest) (*api.CreateRecordResponse, error) {
	return domain_svc.Record().CreateRecord(ctx, req)
}

// UpdateRecord 更新记录
func (r *Record) UpdateRecord(ctx context.Context, req *api.UpdateRecordRequest) (*api.UpdateRecordResponse, error) {
	return domain_svc.Record().UpdateRecord(ctx, req)
}

// DeleteRecord 删除记录
func (r *Record) DeleteRecord(ctx context.Context, req *api.DeleteRecordRequest) (*api.DeleteRecordResponse, error) {
	return domain_svc.Record().DeleteRecord(ctx, req)
}
