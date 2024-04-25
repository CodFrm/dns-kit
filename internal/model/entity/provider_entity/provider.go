package provider_entity

type Platform string

const (
	PlatformCloudflare Platform = "cloudflare"
	PlatformDnsPod     Platform = "dnspod"
)

type Provider struct {
	ID         int64    `gorm:"column:id;type:bigint(20);not null;primary_key"`
	Name       string   `gorm:"column:name;type:varchar(128);not null"`     // 名称
	UserID     string   `gorm:"column:user_id;type:varchar(128);not null"`  // 平台用户id
	Secret     string   `gorm:"column:secret;type:varchar(256);not null"`   // 密钥信息(JSON格式保存)
	Platform   Platform `gorm:"column:platform;type:varchar(128);not null"` // 平台
	Status     int32    `gorm:"column:status;type:tinyint(4);default:0;not null"`
	Createtime int64    `gorm:"column:createtime;type:bigint(20)"`
	Updatetime int64    `gorm:"column:updatetime;type:bigint(20)"`
}
