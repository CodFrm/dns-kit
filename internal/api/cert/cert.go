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
	Expiretime int64                  `json:"expiretime"`
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
	mux.Meta  `path:"/cert" method:"POST"`
	Email     string   `json:"email" binding:"required,email"`
	Domains   []string `json:"domains"` // 域名
	ignoreMsg bool     // 忽略消息推送
}

func (c *CreateRequest) SetIgnoreMsg(ignoreMsg bool) *CreateRequest {
	c.ignoreMsg = ignoreMsg
	return c
}

func (c *CreateRequest) GetIgnoreMsg() bool {
	return c.ignoreMsg
}

// CreateResponse 创建证书响应
type CreateResponse struct {
	ID int64 `json:"id"`
}

// DeleteRequest 删除证书
type DeleteRequest struct {
	mux.Meta `path:"/cert/:id" method:"DELETE"`
	ID       int64 `uri:"id" binding:"required"`
}

// DeleteResponse 删除证书响应
type DeleteResponse struct {
}

// DownloadRequest 下载证书
type DownloadRequest struct {
	mux.Meta `path:"/cert/:id/download" method:"GET"`
	ID       int64 `uri:"id" binding:"required"`
	Type     int   `form:"type"`
}

// DownloadResponse 下载证书响应
type DownloadResponse struct {
	CSR  string `json:"csr"`
	Cert string `json:"data"`
	Key  string `json:"key"`
}
