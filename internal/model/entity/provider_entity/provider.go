package provider_entity

import (
	"context"
	"encoding/json"
	"github.com/codfrm/cago/pkg/i18n"
	"github.com/codfrm/dns-kit/internal/pkg/code"
	"github.com/codfrm/dns-kit/pkg/dns"
	"github.com/codfrm/dns-kit/pkg/dns/provider/cloudflare"
	"github.com/codfrm/dns-kit/pkg/dns/provider/dnspod"
)

type Platform string

const (
	PlatformCloudflare Platform = "cloudflare"
	PlatformTencent    Platform = "tencent"
)

type Provider struct {
	ID         int64  `gorm:"column:id;type:bigint(20);not null;primary_key"`
	Name       string `gorm:"column:name;type:varchar(128);not null"`    // 名称
	UserID     string `gorm:"column:user_id;type:varchar(128);not null"` // 平台用户id
	Secret     string `gorm:"column:secret;type:varchar(256);not null"`  // 密钥信息(JSON格式保存)
	secretMap  map[string]string
	Platform   Platform `gorm:"column:platform;type:varchar(128);not null"` // 平台
	Status     int32    `gorm:"column:status;type:tinyint(4);default:0;not null"`
	Createtime int64    `gorm:"column:createtime;type:bigint(20)"`
	Updatetime int64    `gorm:"column:updatetime;type:bigint(20)"`
}

func (p *Provider) Check(ctx context.Context) error {
	if p == nil {
		return i18n.NewNotFoundError(ctx, code.ProviderNotFound)
	}
	return nil
}

func (p *Provider) SecretMap() map[string]string {
	if p.secretMap != nil {
		return p.secretMap
	}
	ret := make(map[string]string)
	if p.Secret != "" {
		_ = json.Unmarshal([]byte(p.Secret), &ret)
	}
	p.secretMap = ret
	return ret
}

func (p *Provider) Factory(ctx context.Context) (dns.DomainManager, error) {
	var (
		manager dns.DomainManager
		err     error
	)
	switch p.Platform {
	case PlatformCloudflare:
		manager, err = cloudflare.NewCloudflare(p.SecretMap()["key"], p.SecretMap()["email"])
	case PlatformTencent:
		manager, err = dnspod.NewDnsPod(p.SecretMap()["secret_id"], p.SecretMap()["secret_key"])
	default:
		return nil, i18n.NewError(ctx, code.ProviderNotSupport)
	}
	if err != nil {
		return nil, err
	}
	return manager, nil
}
