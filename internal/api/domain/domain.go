package domain

import (
	"github.com/codfrm/cago/pkg/utils/httputils"
	"github.com/codfrm/cago/server/mux"
)

type Item struct {
	ID           int64  `json:"id"`
	Domain       string `json:"domain"`
	ProviderName string `json:"provider_name"`
}

// ListRequest 获取域名列表
type ListRequest struct {
	mux.Meta              `path:"/domain" method:"GET"`
	httputils.PageRequest `form:",inline"`
}

type ListResponse struct {
	httputils.PageResponse[*Item] `json:",inline"`
}

type QueryItem struct {
	ProviderID   int64  `json:"provider_id"`
	ProviderName string `json:"provider_name"`
	DomainID     string `json:"domain_id"`
	Domain       string `json:"domain"`
	IsManaged    bool   `json:"is_managed"`
}

// QueryRequest 查询域名列表
type QueryRequest struct {
	mux.Meta `path:"/domain/query" method:"GET"`
}

type QueryResponse struct {
	Items []*QueryItem `json:"items"`
}

// AddRequest 纳管域名
type AddRequest struct {
	mux.Meta   `path:"/domain" method:"POST"`
	ProviderID int64  `json:"provider_id"`
	DomainID   string `json:"domain_id"`
	Domain     string `json:"domain"`
}

type AddResponse struct {
}

// DeleteRequest 删除域名
type DeleteRequest struct {
	mux.Meta `path:"/domain/:id" method:"DELETE"`
	ID       int64 `uri:"id"`
}

type DeleteResponse struct {
}
