package code

import "github.com/codfrm/cago/pkg/i18n"

func init() {
	i18n.Register(i18n.DefaultLang, zhCN)
}

var zhCN = map[int]string{
	UserIsBanned:          "用户已被禁用",
	UserNotFound:          "用户不存在",
	UserNotLogin:          "用户未登录",
	UsernameAlreadyExists: "用户名已存在",

	ProviderNotSupport:  "供应商不支持",
	ProviderSecretError: "供应商密钥错误",
	ProviderExist:       "供应商已存在",
	ProviderNotFound:    "供应商不存在",
}
