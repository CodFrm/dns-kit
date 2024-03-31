package user

import (
	"github.com/codfrm/cago/server/mux"
)

// LoginRequest 登录
type LoginRequest struct {
	mux.Meta `path:"/user/login" method:"POST"`
	// 用户名
	Username string `form:"username" binding:"required"`
	// 密码
	Password string `form:"password" binding:"required"`
}

type LoginResponse struct {
}

// LogoutRequest 登出
type LogoutRequest struct {
	mux.Meta `path:"/user/logout" method:"GET"`
}

type LogoutResponse struct {
}

// CurrentUserRequest 获取当前用户
type CurrentUserRequest struct {
	mux.Meta `path:"/user/current" method:"GET"`
}

type CurrentUserResponse struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
}
