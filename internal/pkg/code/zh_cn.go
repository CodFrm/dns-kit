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

	DomainNotFound:  "域名不存在或者不在纳管中",
	DomainIsManaged: "域名已纳管",

	RecordCreateFailed: "记录创建失败: %v",
	RecordUpdateFailed: "记录更新失败: %v",

	InvalidDomain:   "无效的域名",
	CertNotFound:    "证书不存在",
	CertNotActive:   "证书未激活",
	CertStatusApply: "证书申请中",

	CDNNotFound: "CDN不存在",
}
