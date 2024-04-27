package domain_svc

import (
	"context"

	api "github.com/codfrm/dns-kit/internal/api/domain"
)

type RecordSvc interface {
	// RecordList 获取记录列表
	RecordList(ctx context.Context, req *api.RecordListRequest) (*api.RecordListResponse, error)
}

type recordSvc struct {
}

var defaultRecord = &recordSvc{}

func Record() RecordSvc {
	return defaultRecord
}

// RecordList 获取记录列表
func (r *recordSvc) RecordList(ctx context.Context, req *api.RecordListRequest) (*api.RecordListResponse, error) {
	return nil, nil
}
