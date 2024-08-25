package project

import "time"

type Version struct {
	ID         uint      `gorm:"column:id"`
	ProjectID  uint      `gorm:"index;not null;column:project_id"`
	Name       string    `gorm:"not null;column:name"`
	StorageKey string    `gorm:"not null;column:storage_key"`
	CreatedAt  time.Time `gorm:"not null;column:created_at"`
}

type Project struct {
	ID               uint      `gorm:"column:id"`
	Name             string    `gorm:"not null;column:name"`
	NumberOfVersions uint      `gorm:"column:number_of_versions"`
	CreatedAt        time.Time `gorm:"not null;column:created_at"`
}
