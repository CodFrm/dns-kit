package migrations

import (
	"github.com/codfrm/dns-kit/internal/model/entity/acme_entity"
	"github.com/codfrm/dns-kit/internal/model/entity/cert_entity"
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func T20240427() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "20240427",
		Migrate: func(tx *gorm.DB) error {
			entities := []any{
				&cert_entity.Cert{},
				&acme_entity.Acme{},
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
