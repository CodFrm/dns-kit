package user_entity

import (
	"context"

	"github.com/codfrm/cago/pkg/consts"
	"github.com/codfrm/cago/pkg/i18n"
	"github.com/codfrm/dns-kit/internal/pkg/code"
)

type User struct {
	ID             int64  `gorm:"column:id;not null;primary_key"`
	Username       string `gorm:"column:username;type:varchar(255);index:username,unique;not null"`
	HashedPassword string `gorm:"column:hashed_password;type:varchar(255);not null"`
	Status         int    `gorm:"column:status;type:int(11);not null"`
	Createtime     int64  `gorm:"column:createtime;type:bigint(20)"`
	Updatetime     int64  `gorm:"column:updatetime;type:bigint(20)"`
}

func (u *User) Check(ctx context.Context) error {
	if u == nil {
		return i18n.NewError(ctx, code.UserNotFound)
	}
	if u.Status != consts.ACTIVE {
		return i18n.NewError(ctx, code.UserIsBanned)
	}
	return nil
}
