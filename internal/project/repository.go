//go:generate mockgen --source=repository.go -destination=repository_mock.go -package=project

package project

import (
	"context"
	"errors"

	"gorm.io/gorm"
)

var (
	ErrNotFound   = errors.New("not found")
	ErrNotDeleted = errors.New("not deleted")
)

type Repository interface {
	GetAllProjects(ctx context.Context) ([]Project, error)
	GetProjectByID(ctx context.Context, id int) (*Project, error)
	GetProjectVersions(ctx context.Context, projectID int) ([]Version, error)
	GetVersionByIDAndProjectID(ctx context.Context, versionID, projectID uint) (*Version, error)
	DeleteVersionByIDTransaction(ctx context.Context, versionID uint) (*gorm.DB, error)
}

type RepositoryImpl struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &RepositoryImpl{
		db: db,
	}
}

func (r *RepositoryImpl) GetAllProjects(ctx context.Context) ([]Project, error) {
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

	return projects, nil
}

func (r *RepositoryImpl) GetProjectByID(ctx context.Context, id int) (*Project, error) {
	var project Project

	result := r.db.WithContext(ctx).
		Select("projects.id, projects.name, COUNT(versions.id) as number_of_versions, projects.created_at").
		Joins("LEFT JOIN versions ON projects.id = versions.project_id").
		Group("projects.id").
		First(&project, id)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}

		return nil, result.Error
	}

	return &project, nil
}

func (r *RepositoryImpl) GetProjectVersions(ctx context.Context, projectID int) ([]Version, error) {
	var versions []Version

	result := r.db.WithContext(ctx).
		Where("project_id = ?", projectID).
		Find(&versions)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}

		return nil, result.Error
	}

	return versions, nil
}

func (r *RepositoryImpl) GetVersionByIDAndProjectID(ctx context.Context, versionID, projectID uint) (*Version, error) {
	var version Version

	result := r.db.WithContext(ctx).
		Where("project_id = ? AND id = ?", projectID, versionID).
		First(&version)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}

		return nil, result.Error
	}

	return &version, nil
}

func (r *RepositoryImpl) DeleteVersionByIDTransaction(ctx context.Context, versionID uint) (*gorm.DB, error) {
	tx := r.db.WithContext(ctx).Begin()
	if err := tx.Error; err != nil {
		return nil, err
	}

	result := tx.Delete(&Version{}, versionID)
	if result.Error != nil {
		tx.Rollback()
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		tx.Rollback()
		return nil, ErrNotDeleted
	}

	return tx, nil
}
