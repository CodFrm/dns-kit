package migrations

import (
	"github.com/codfrm/dns-kit/internal/model/entity/domain_entity"
	"github.com/codfrm/dns-kit/internal/model/entity/provider_entity"
	"github.com/codfrm/dns-kit/internal/model/entity/user_entity"
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func T20240427() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "20240427",
		Migrate: func(tx *gorm.DB) error {
			// 初始化用户
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
			return nil
		},
		Rollback: func(tx *gorm.DB) error {
			return nil
		},
	}
}
