package project

import (
	"errors"

	"github.com/uniwise/parrot/internal/cache"
	"github.com/uniwise/parrot/internal/poedit"
)

type Service interface {
	GetTranslation(projectID int, languageCode, format string) (data []byte, err error)
	PurgeTranslation(projectID int, languageCode string) (err error)
	PurgeProject(projectID int) (err error)
}

type ServiceImpl struct {
	Client poedit.Client
	Cache  cache.Cache
}

func NewService(cli poedit.Client, cache cache.Cache) *ServiceImpl {
	return &ServiceImpl{
		Client: cli,
		Cache:  cache,
	}
}

func (s *ServiceImpl) GetTranslation(projectID int, languageCode, format string) ([]byte, error) {
	return nil, errors.New("Not implemented")
}

func (s *ServiceImpl) PurgeTranslation(projectID int, languageCode string) error {
	return errors.New("Not implemented")
}

func (s *ServiceImpl) PurgeProject(projectID int) error {
	return errors.New("Not implemented")
}
