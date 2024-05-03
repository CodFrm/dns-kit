package cert

import (
	"github.com/codfrm/cago/pkg/utils/httputils"
	"github.com/codfrm/cago/server/mux"
	"github.com/codfrm/dns-kit/internal/model/entity/cert_entity"
)

type Item struct {
	ID         int64                  `json:"id"`
	Email      string                 `json:"email"`
	Domains    []string               `json:"domains"`
	Status     cert_entity.CertStatus `json:"status"`
	Createtime int64                  `json:"createtime"`
}

// ListRequest 获取证书列表
type ListRequest struct {
	mux.Meta              `path:"/cert" method:"GET"`
	httputils.PageRequest `json:",inline"`
}

// ListResponse 证书列表
type ListResponse struct {
	httputils.PageResponse[*Item] `json:",inline"`
}

// CreateRequest 创建证书
type CreateRequest struct {
	mux.Meta `path:"/cert" method:"POST"`
	Email    string   `json:"email"`
	Domains  []string `json:"domains"` // 域名
}

// CreateResponse 创建证书响应
type CreateResponse struct {
}
