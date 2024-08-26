package migrations

import (
	"time"

	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

var Migration00002AddStorageKeyToVersionsTable = &gormigrate.Migration{
	ID: "00002_add_storage_key_to_versions_table",
	Migrate: func(tx *gorm.DB) error {
		type Version struct {
			ID         uint   `gorm:"primaryKey"`
			ProjectID  uint   `gorm:"index:idx_name_projectId,unique;not null"`
			Name       string `gorm:"index:idx_name_projectId,unique;not null;size:191"`
			StorageKey string `gorm:"not null;size:191"`
			CreatedAt  time.Time
		}

		type Project struct {
			ID              uint      `gorm:"primaryKey"`
			ClientProjectID uint      `gorm:"not null"` // POeditor Project ID
			Name            string    `gorm:"not null;size:191"`
			Versions        []Version `gorm:"constraint:OnDelete:CASCADE;"`
			CreatedAt       time.Time
		}

		if err := tx.Migrator().DropConstraint(&Project{}, "Versions"); err != nil {
			return err
		}

		return tx.AutoMigrate(&Project{}, &Version{})
	},
}
