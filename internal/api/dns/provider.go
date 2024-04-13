package dns

import (
	"github.com/codfrm/cago/pkg/utils/httputils"
	"github.com/codfrm/cago/server/mux"
	"github.com/codfrm/dns-kit/internal/model/entity/dns_provider_entity"
)

type ProviderItem struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

// ListProviderRequest 获取dns提供商列表
type ListProviderRequest struct {
	mux.Meta              `path:"/dns/provider" method:"GET"`
	httputils.PageRequest `form:",inline"`
}

type ListProviderResponse struct {
	httputils.PageResponse[*ProviderItem] `json:",inline"`
}

// CreateProviderRequest 创建dns提供商
type CreateProviderRequest struct {
	mux.Meta `path:"/dns/provider" method:"POST"`
	Name     string                          `json:"name" binding:"required,max=128" label:"名称"`
	Platform dns_provider_entity.DnsPlatform `json:"platform" binding:"required,max=128" label:"平台"`
	Secret   map[string]string               `json:"secret" binding:"required,max=256" label:"密钥信息"`
}

type CreateProviderResponse struct {
	ID int64 `json:"id"`
}
