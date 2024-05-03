package cert

import (
	"github.com/codfrm/cago/pkg/utils/httputils"
	"github.com/codfrm/cago/server/mux"
	"github.com/codfrm/dns-kit/internal/model/entity/cert_hosting_entity"
)

type HostingItem struct {
	ID         int64                                 `json:"id"`
	CdnID      int64                                 `json:"cdn_id"`
	CDN        string                                `json:"cdn"`
	CertID     int64                                 `json:"cert_id"`
	Status     cert_hosting_entity.CertHostingStatus `json:"status"`
	Createtime int64                                 `json:"createtime"`
}

// HostingListRequest 托管列表
type HostingListRequest struct {
	mux.Meta              `path:"/cert/hosting" method:"GET"`
	httputils.PageRequest `form:",inline"`
}

// HostingListResponse 托管列表
type HostingListResponse struct {
	httputils.PageResponse[*HostingItem] `json:",inline"`
}

type HostingQueryItem struct {
	ID        int64  `json:"id"`
	Domain    string `json:"domain"`
	IsManaged bool   `json:"is_managed"`
}

// HostingQueryRequest 查询托管
type HostingQueryRequest struct {
	mux.Meta `path:"/cert/hosting/query" method:"GET"`
}

// HostingQueryResponse 查询托管
type HostingQueryResponse struct {
	List []*HostingQueryItem `json:"list"`
}

// HostingAddRequest 添加托管
type HostingAddRequest struct {
	mux.Meta `path:"/cert/hosting" method:"POST"`
	Email    string `json:"email" binding:"required,email"`
	CdnID    int64  `json:"cdn_id"`
}

// HostingAddResponse 添加托管
type HostingAddResponse struct {
}

// HostingDeleteRequest 删除托管
type HostingDeleteRequest struct {
	mux.Meta `path:"/cert/hosting/:id" method:"DELETE"`
	ID       int64 `uri:"id"`
}

// HostingDeleteResponse 删除托管
type HostingDeleteResponse struct {
}
