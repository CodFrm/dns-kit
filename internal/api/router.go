package api

import (
	"context"

	_ "github.com/codfrm/cago/examples/simple/docs"
	"github.com/codfrm/cago/server/mux"
)

// Router 路由
// @title    api文档
// @version  1.0
// @BasePath /api/v1
func Router(ctx context.Context, root *mux.Router) error {
	// 注册储存实例
	_ = root.Group("/api/v1")

	return nil
}
