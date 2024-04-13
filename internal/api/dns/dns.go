package dns

import (
	"github.com/codfrm/cago/pkg/utils/httputils"
	"github.com/codfrm/cago/server/mux"
)

type Item struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	Platform string `json:"platform"`
}

// ListRequest 获取dns列表
type ListRequest struct {
	mux.Meta              `path:"/dns" method:"GET"`
	httputils.PageRequest `form:",inline"`
}

type ListResponse struct {
	httputils.PageResponse[*Item] `json:",inline"`
}
