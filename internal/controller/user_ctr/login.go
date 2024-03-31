package user_ctr

import (
	"context"

	api "github.com/codfrm/dns-kit/internal/api/user"
	"github.com/codfrm/dns-kit/internal/service/user_svc"
)

type Login struct {
}

func NewLogin() *Login {
	return &Login{}
}

// Login 登录
func (l *Login) Login(ctx context.Context, req *api.LoginRequest) (*api.LoginResponse, error) {
	return user_svc.Login().Login(ctx, req)
}

// Logout 登出
func (l *Login) Logout(ctx context.Context, req *api.LogoutRequest) (*api.LogoutResponse, error) {
	return user_svc.Login().Logout(ctx, req)
}

// CurrentUser 获取当前用户
func (l *Login) CurrentUser(ctx context.Context, req *api.CurrentUserRequest) (*api.CurrentUserResponse, error) {
	return user_svc.Login().CurrentUser(ctx, req)
}
