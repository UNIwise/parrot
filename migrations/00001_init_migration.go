package migrations

import (
	"time"

	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

// InitMigration is the first migration that creates the initial tables
var Migration00001Init = &gormigrate.Migration{
	ID: "00001_init",
	Migrate: func(tx *gorm.DB) error {
		type Version struct {
			ID        uint `gorm:"primaryKey"`
			ProjectID uint
			Name      string
			CreatedAt time.Time
		}

		type Project struct {
			ID        uint `gorm:"primaryKey"`
			Name      string
			Versions  []Version `gorm:"constraint:OnDelete:CASCADE;"`
			CreatedAt time.Time
		}

		return tx.AutoMigrate(&Project{}, &Version{})
	},
}
