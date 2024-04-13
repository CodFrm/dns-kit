package dns_provider_entity

type DnsPlatform string

const (
	DnsPlatformCloudflare DnsPlatform = "cloudflare"
	DnsPlatformDnsPod     DnsPlatform = "dnspod"
)

type DnsProvider struct {
	ID         int64       `gorm:"column:id;type:bigint(20);not null;primary_key"`
	Name       string      `gorm:"column:name;type:varchar(128);not null"`     // 名称
	UserID     string      `gorm:"column:user_id;type:varchar(128);not null"`  // 平台用户id
	Secret     string      `gorm:"column:secret;type:varchar(256);not null"`   // 密钥信息(JSON格式保存)
	Platform   DnsPlatform `gorm:"column:platform;type:varchar(128);not null"` // 平台
	Status     int32       `gorm:"column:status;type:tinyint(4);default:0;not null"`
	Createtime int64       `gorm:"column:createtime;type:bigint(20)"`
	Updatetime int64       `gorm:"column:updatetime;type:bigint(20)"`
}
