package domain

import (
	"github.com/codfrm/cago/server/mux"
	"github.com/codfrm/dns-kit/pkg/dns"
)

// RecordListRequest 获取记录列表
type RecordListRequest struct {
	mux.Meta `path:"/domain/:id/record" method:"GET"`
	ID       int64 `uri:"id"`
}

type RecordListResponse struct {
	List  []*dns.Record `json:"list"`
	Extra []*dns.Extra  `json:"extra"`
}

// CreateRecordRequest 创建记录
type CreateRecordRequest struct {
	mux.Meta `path:"/domain/:id/record" method:"POST"`
	ID       int64          `uri:"id"`
	Type     dns.RecordType `json:"type" binding:"required,oneof=A AAAA CNAME TXT MX NS"`
	Name     string         `json:"name" binding:"required"`
	Value    string         `json:"value" binding:"required"`
	TTL      int            `json:"ttl"`
	Extra    map[string]any `json:"extra"`
}

type CreateRecordResponse struct {
}

// UpdateRecordRequest 更新记录
type UpdateRecordRequest struct {
	mux.Meta `path:"/domain/:id/record/:recordID" method:"PUT"`
	ID       int64          `uri:"id"`
	RecordID string         `uri:"recordID"`
	Type     dns.RecordType `json:"type" binding:"required,oneof=A AAAA CNAME TXT MX NS"`
	Name     string         `json:"name" binding:"required"`
	Value    string         `json:"value" binding:"required"`
	TTL      int            `form:"ttl,default=600" json:"ttl" binding:"required"`
	Extra    map[string]any `json:"extra"`
}

type UpdateRecordResponse struct {
}

// DeleteRecordRequest 删除记录
type DeleteRecordRequest struct {
	mux.Meta `path:"/domain/:id/record/:recordID" method:"DELETE"`
	ID       int64  `uri:"id"`
	RecordID string `uri:"recordID"`
}

type DeleteRecordResponse struct {
}
