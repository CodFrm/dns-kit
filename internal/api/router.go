package api

import (
	"context"

	"github.com/codfrm/cago/server/mux"
	_ "github.com/codfrm/dns-kit/docs"
	"github.com/codfrm/dns-kit/internal/controller/domain_ctr"
	"github.com/codfrm/dns-kit/internal/controller/provider_ctr"
	"github.com/codfrm/dns-kit/internal/controller/user_ctr"
	"github.com/codfrm/dns-kit/internal/service/user_svc"
)

// Router 路由
// @title    api文档
// @version  1.0
// @BasePath /api/v1
func Router(ctx context.Context, root *mux.Router) error {
	r := root.Group("/api/v1")

	userCtr := user_ctr.NewUser()
	{
		// 绑定路由
		r.Group("/", user_svc.User().AuditMiddleware("user")).Bind(
			userCtr.Register,
			userCtr.Login,
		)

		r.Group("/", user_svc.User().Middleware(true)).Bind(
			userCtr.CurrentUser,
			userCtr.Logout,
			userCtr.RefreshToken,
		)
	}

	domainCtr := domain_ctr.NewDomain()
	{
		r.Group("/", user_svc.User().Middleware(true)).Bind(
			domainCtr.List,
			domainCtr.Query,
		)

		r.Group("/", user_svc.User().Middleware(true),
			user_svc.User().AuditMiddleware("domain")).Bind(
			domainCtr.Add,
			domainCtr.Delete,
		)
	}

	providerCtr := provider_ctr.NewProvider()
	{
		r.Group("/", user_svc.User().Middleware(true)).Bind(
			providerCtr.ListProvider,
		)

		r.Group("/", user_svc.User().Middleware(true),
			user_svc.User().AuditMiddleware("provider")).Bind(
			providerCtr.CreateProvider,
			providerCtr.UpdateProvider,
			providerCtr.DeleteProvider,
		)
	}

	return nil
}
