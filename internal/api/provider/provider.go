package provider

import (
	"github.com/codfrm/cago/pkg/utils/httputils"
	"github.com/codfrm/cago/server/mux"
	"github.com/codfrm/dns-kit/internal/model/entity/provider_entity"
)

type Item struct {
	ID       int64                    `json:"id"`
	Name     string                   `json:"name"`
	Platform provider_entity.Platform `json:"platform"`
}

// ListProviderRequest 获取供应商列表
type ListProviderRequest struct {
	mux.Meta              `path:"/provider" method:"GET"`
	httputils.PageRequest `form:",inline"`
}

type ListProviderResponse struct {
	httputils.PageResponse[*Item] `json:",inline"`
}

// CreateProviderRequest 创建供应商
type CreateProviderRequest struct {
	mux.Meta `path:"/provider" method:"POST"`
	Name     string                   `json:"name" binding:"required,max=128" label:"名称"`
	Platform provider_entity.Platform `json:"platform" binding:"required,max=128" label:"平台"`
	Secret   map[string]string        `json:"secret" binding:"required,max=256" label:"密钥信息"`
}

type CreateProviderResponse struct {
}

// UpdateProviderRequest 更新供应商
type UpdateProviderRequest struct {
	mux.Meta `path:"/provider/:id" method:"PUT"`
	ID       int64             `uri:"id" binding:"required" label:"ID"`
	Name     string            `json:"name" binding:"required,max=128" label:"名称"`
	Secret   map[string]string `json:"secret" binding:"required,max=256" label:"密钥信息"`
}

type UpdateProviderResponse struct {
}

// DeleteProviderRequest 删除供应商
type DeleteProviderRequest struct {
	mux.Meta `path:"/provider/:id" method:"DELETE"`
	ID       int64 `uri:"id" binding:"required" label:"ID"`
}

type DeleteProviderResponse struct {
}
