package migrations

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func T20240326() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "20240326",
		Migrate: func(tx *gorm.DB) error {
			return nil
		},
		Rollback: func(tx *gorm.DB) error {
			return nil
		},
	}
}
