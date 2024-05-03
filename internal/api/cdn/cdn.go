package cdn

import (
	"github.com/codfrm/cago/pkg/utils/httputils"
	"github.com/codfrm/cago/server/mux"
)

type Item struct {
	ID           int64  `json:"id"`
	ProviderName string `json:"provider_name"`
	Domain       string `json:"domain"`
	Createtime   int64  `json:"createtime"`
}

// ListRequest 获取纳管的cdn列表
type ListRequest struct {
	mux.Meta              `path:"/cdn" method:"GET"`
	httputils.PageRequest `form:",inline"`
}

type ListResponse struct {
	httputils.PageResponse[*Item] `json:",inline"`
}

type QueryItem struct {
	ProviderID   int64  `json:"provider_id"`
	ProviderName string `json:"provider_name"`
	ID           string `json:"id"`
	Domain       string `json:"domain"`
	IsManaged    bool   `json:"is_managed"`
}

// QueryRequest 查询cdn
type QueryRequest struct {
	mux.Meta `path:"/cdn/query" method:"GET"`
}

type QueryResponse struct {
	Items []*QueryItem `json:"items"`
}

// AddRequest 添加cdn进入纳管
type AddRequest struct {
	mux.Meta   `path:"/cdn" method:"POST"`
	ProviderID int64  `json:"provider_id"`
	ID         string `json:"id"`
	Domain     string `json:"domain"`
}

type AddResponse struct {
}

// DeleteRequest 删除cdn
type DeleteRequest struct {
	mux.Meta `path:"/cdn/:id" method:"DELETE"`
	ID       int64 `uri:"id" binding:"required"`
}

type DeleteResponse struct {
}
