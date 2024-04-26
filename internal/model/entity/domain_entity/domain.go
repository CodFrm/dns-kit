package domain_entity

type Domain struct {
	ID         int64  `gorm:"column:id;type:bigint(20);not null;primary_key"`
	ProviderID int64  `gorm:"column:provider_id;type:bigint(20);not null"` // 供应商id
	Name       string `gorm:"column:name;type:varchar(128);not null"`      // 名称
	Domain     string `gorm:"column:domain;type:varchar(128);not null"`    // 域名
	Status     int32  `gorm:"column:status;type:tinyint(4);default:0;not null"`
	Createtime int64  `gorm:"column:createtime;type:bigint(20)"`
	Updatetime int64  `gorm:"column:updatetime;type:bigint(20)"`
}
