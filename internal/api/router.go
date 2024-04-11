package api

import (
	"context"
	"github.com/codfrm/cago/pkg/iam/authn"
	"github.com/codfrm/cago/server/mux"
	_ "github.com/codfrm/dns-kit/docs"
	"github.com/codfrm/dns-kit/internal/controller/user_ctr"
	"github.com/codfrm/dns-kit/internal/repository/user_repo"
	"github.com/codfrm/dns-kit/internal/service/user_svc"
)

// Router 路由
// @title    api文档
// @version  1.0
// @BasePath /api/v1
func Router(ctx context.Context, root *mux.Router) error {
	// 注册认证模块
	auth := authn.New(user_repo.User(),
		authn.WithMiddleware(user_svc.User().Middleware()),
	)
	authn.SetDefault(auth)

	r := root.Group("/api/v1")

	userLoginCtr := user_ctr.NewUser()
	{
		// 绑定路由
		r.Group("/").Bind(
			userLoginCtr.Register,
			userLoginCtr.Login,
		)

		r.Group("/", auth.Middleware(true)).Bind(
			userLoginCtr.CurrentUser,
			userLoginCtr.Logout,
			userLoginCtr.RefreshToken,
		)
	}

	return nil
}
