package domain

import "github.com/codfrm/cago/server/mux"

type Record struct {
}

// RecordListRequest 获取记录列表
type RecordListRequest struct {
	mux.Meta `path:"/domain/:id/record" method:"GET"`
	ID       int64 `uri:"id"`
}

type RecordListResponse struct {
	List []*Record
}
