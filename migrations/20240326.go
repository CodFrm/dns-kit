package migrations

import (
	"context"
	"github.com/codfrm/cago/database/db"
	"github.com/codfrm/dns-kit/internal/api/user"
	"github.com/codfrm/dns-kit/internal/model/entity/user_entity"
	"github.com/codfrm/dns-kit/internal/service/user_svc"
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func T20240326() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "20240326",
		Migrate: func(tx *gorm.DB) error {
			// 初始化用户
			ctx := context.Background()
			ctx = db.ContextWithDB(ctx, tx)
			if err := tx.Migrator().AutoMigrate(&user_entity.User{}); err != nil {
				return err
			}
			// 添加admin用户
			_, err := user_svc.Login().Register(ctx, &user.RegisterRequest{
				Username: "admin",
				Password: "123456",
			})
			if err != nil {
				return err
			}
			return nil
		},
		Rollback: func(tx *gorm.DB) error {
			return nil
		},
	}
}
