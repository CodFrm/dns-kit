package user_svc

import (
	"context"

	api "github.com/codfrm/dns-kit/internal/api/user"
)

type LoginSvc interface {
	// Login 登录
	Login(ctx context.Context, req *api.LoginRequest) (*api.LoginResponse, error)
	// Logout 登出
	Logout(ctx context.Context, req *api.LogoutRequest) (*api.LogoutResponse, error)
	// CurrentUser 获取当前用户
	CurrentUser(ctx context.Context, req *api.CurrentUserRequest) (*api.CurrentUserResponse, error)
}

type loginSvc struct {
}

var defaultLogin = &loginSvc{}

func Login() LoginSvc {
	return defaultLogin
}

// Login 登录
func (l *loginSvc) Login(ctx context.Context, req *api.LoginRequest) (*api.LoginResponse, error) {
	return nil, nil
}

// Logout 登出
func (l *loginSvc) Logout(ctx context.Context, req *api.LogoutRequest) (*api.LogoutResponse, error) {
	return nil, nil
}

// CurrentUser 获取当前用户
func (l *loginSvc) CurrentUser(ctx context.Context, req *api.CurrentUserRequest) (*api.CurrentUserResponse, error) {
	return nil, nil
}
