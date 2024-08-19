package project

import (
	"gorm.io/gorm"
)

type Repository interface {
	// TODO: Add methods to implement
}

type RepositoryImpl struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &RepositoryImpl{
		db: db,
	}
}