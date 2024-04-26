package migrations

import (
	"context"

	"github.com/codfrm/cago/database/db"
	"github.com/codfrm/dns-kit/internal/api/user"
	"github.com/codfrm/dns-kit/internal/model/entity/domain_entity"
	"github.com/codfrm/dns-kit/internal/model/entity/provider_entity"
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
			ctx = db.WithContextDB(ctx, tx)
			entities := []any{
				&user_entity.User{},
				&domain_entity.Domain{},
				&provider_entity.Provider{},
			}
			for _, entity := range entities {
				if err := tx.Migrator().AutoMigrate(entity); err != nil {
					return err
				}
			}
			// 添加admin用户
			_, err := user_svc.User().Register(ctx, &user.RegisterRequest{
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
