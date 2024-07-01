package migrations

import (
	"github.com/codfrm/dns-kit/internal/model/entity/cert_hosting_entity"
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func T20240701() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "20240701",
		Migrate: func(tx *gorm.DB) error {
			entities := []any{
				&cert_hosting_entity.CertHosting{},
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
