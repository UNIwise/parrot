package project

import (
	"context"

	"io/ioutil"
	"net/http"

	"github.com/pkg/errors"
	"github.com/uniwise/parrot/internal/cache"
	"github.com/uniwise/parrot/internal/poedit"
)

type Service interface {
	GetTranslation(ctx context.Context, projectID int, languageCode, format string) (data []byte, err error)
	PurgeTranslation(ctx context.Context, projectID int, languageCode string) (err error)
	PurgeProject(ctx context.Context, projectID int) (err error)
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

func (s *ServiceImpl) GetTranslation(ctx context.Context, projectID int, languageCode, format string) ([]byte, error) {
	data, err := s.Cache.GetTranslation(ctx, projectID, languageCode, format)
	if err != nil && !errors.Is(err, cache.ErrCacheMiss) {
		return nil, err
	}
	if err == nil {
		return data, nil
	}

	resp, err := s.Client.ExportProject(ctx, poedit.ExportProjectRequest{
		ID:       projectID,
		Language: languageCode,
		Type:     format,
		Filters:  []string{"translated"},
	})
	if err != nil {
		return nil, err
	}

	// TODO: Make use of injected http client, to support timeouts
	d, err := http.Get(resp.Result.URL)
	if err != nil {
		return nil, err
	}
	defer d.Body.Close()

	if d.StatusCode != http.StatusOK {
		return nil, errors.Errorf("Response code '%d' from download GET", d.StatusCode)
	}

	data, err = ioutil.ReadAll(d.Body)
	if err != nil {
		return nil, err
	}

	if err := s.Cache.SetTranslation(ctx, projectID, languageCode, format, data); err != nil {
		return nil, err
	}

	return data, nil
}

func (s *ServiceImpl) PurgeTranslation(ctx context.Context, projectID int, languageCode string) error {
	return errors.New("Not implemented")
}

func (s *ServiceImpl) PurgeProject(ctx context.Context, projectID int) error {
	return errors.New("Not implemented")
}
