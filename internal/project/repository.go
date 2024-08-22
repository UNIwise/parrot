//go:generate mockgen --source=repository.go -destination=repository_mock.go -package=project

package project

import (
	"context"
	"errors"

	"gorm.io/gorm"
)

var (
	ErrNotFound = errors.New("not found")
)

type Repository interface {
	GetAllProjects(ctx context.Context) (*[]Project, error)
}

type RepositoryImpl struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &RepositoryImpl{
		db: db,
	}
}

func (r *RepositoryImpl) GetAllProjects(ctx context.Context) (*[]Project, error) {
	var projects []Project

	result := r.db.WithContext(ctx).
		Select("projects.id, projects.name, COUNT(versions.id) as number_of_versions, projects.created_at").
		Joins("LEFT JOIN versions ON projects.id = versions.project_id").
		Group("projects.id").
		Find(&projects)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}

		return nil, result.Error
	}

	return &projects, nil
}
