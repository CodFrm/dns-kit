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
