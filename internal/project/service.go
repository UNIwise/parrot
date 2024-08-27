//go:generate mockgen --source=service.go -destination=service_mock.go -package=project

package project

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	gosundheit "github.com/AppsFlyer/go-sundheit"
	"github.com/AppsFlyer/go-sundheit/checks"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/uniwise/parrot/internal/cache"
	"github.com/uniwise/parrot/internal/storage"
	"github.com/uniwise/parrot/pkg/poedit"
	"golang.org/x/sync/semaphore"
)

type Translation struct {
	TTL      time.Duration
	Checksum string
	Data     []byte
}

type Service interface {
	GetTranslation(ctx context.Context, projectID int, languageCode, format string) (trans *Translation, err error)
	PurgeTranslation(ctx context.Context, projectID int, languageCode string) (err error)
	PurgeProject(ctx context.Context, projectID int) (err error)
	RegisterChecks(h gosundheit.Health) (err error)
	GetAllProjects(ctx context.Context) ([]Project, error)
	GetProjectByID(ctx context.Context, id int) (*Project, error)
	GetProjectVersions(ctx context.Context, projectID int) ([]Version, error)
	DeleteProjectVersionByIDAndProjectID(ctx context.Context, ID, projectID uint) error
	CreateLanguagesVersion(ctx context.Context, projectID int, name string) error
}

type ServiceImpl struct {
	Logger            *logrus.Entry
	Client            poedit.Client
	Cache             cache.Cache
	RenewalThreshold  time.Duration
	PreFetchSemaphore *semaphore.Weighted
	storage           storage.Storage
	repo              Repository
}

func NewService(cli poedit.Client, storage storage.Storage, repo Repository, cache cache.Cache, renewalThreshold time.Duration, entry *logrus.Entry) *ServiceImpl {
	return &ServiceImpl{
		Logger:            entry,
		Client:            cli,
		Cache:             cache,
		RenewalThreshold:  renewalThreshold,
		PreFetchSemaphore: semaphore.NewWeighted(1),
		storage:           storage,
		repo:              repo,
	}
}

func (s *ServiceImpl) GetTranslation(ctx context.Context, projectID int, languageCode, format string) (*Translation, error) {
	item, err := s.Cache.GetTranslation(ctx, projectID, languageCode, format)
	if err != nil && !errors.Is(err, cache.ErrCacheMiss) {
		return nil, err
	}
	if err == nil {
		expiresAt := item.CreatedAt.Add(s.Cache.GetTTL())

		if time.Until(expiresAt) < s.RenewalThreshold {
			go func() {
				if !s.PreFetchSemaphore.TryAcquire(1) {
					return
				}

				defer s.PreFetchSemaphore.Release(1)

				s.Logger.Debugf("Pre-fetching language %s format %s for project %d", languageCode, format, projectID)

				_, _, err := s.fetchAndCacheTranslation(context.Background(), projectID, languageCode, format)
				if err != nil {
					s.Logger.Errorf("Failed to pre-fetch language %s format %s for project %d", languageCode, format, projectID)
				}
			}()
		}

		return &Translation{
			TTL:      s.Cache.GetTTL(),
			Checksum: item.Checksum,
			Data:     item.Data,
		}, nil
	}

	data, checksum, err := s.fetchAndCacheTranslation(ctx, projectID, languageCode, format)
	if err != nil {
		return nil, err
	}

	return &Translation{
		TTL:      s.Cache.GetTTL(),
		Checksum: checksum,
		Data:     data,
	}, nil
}

func (s *ServiceImpl) PurgeTranslation(ctx context.Context, projectID int, languageCode string) error {
	return errors.New("Not implemented")
}

func (s *ServiceImpl) PurgeProject(ctx context.Context, projectID int) error {
	return errors.New("Not implemented")
}

func (s *ServiceImpl) fetchAndCacheTranslation(ctx context.Context, projectID int, languageCode, format string) ([]byte, string, error) {
	resp, err := s.Client.ExportProject(ctx, poedit.ExportProjectRequest{
		ID:       projectID,
		Language: languageCode,
		Type:     format,
		Filters:  []string{"translated"},
	})
	if err != nil {
		return nil, "", err
	}

	// TODO: Make use of injected http client, to support timeouts
	d, err := http.Get(resp.Result.URL)
	if err != nil {
		return nil, "", err
	}
	defer d.Body.Close()

	if d.StatusCode != http.StatusOK {
		return nil, "", errors.Errorf("Response code '%d' from download GET", d.StatusCode)
	}

	data, err := ioutil.ReadAll(d.Body)
	if err != nil {
		return nil, "", err
	}

	checksum, err := s.Cache.SetTranslation(ctx, projectID, languageCode, format, data)
	if err != nil {
		return nil, "", err
	}

	return data, checksum, nil
}

func (s *ServiceImpl) RegisterChecks(h gosundheit.Health) error {
	c, err := checks.NewPingCheck("cache", s.Cache)
	if err != nil {
		return errors.Wrap(err, "Failed to instantiate cache healthcheck")
	}

	if err := h.RegisterCheck(c, gosundheit.ExecutionPeriod(time.Second*10)); err != nil {
		return errors.Wrap(err, "Failed to register cache healthcheck")
	}

	return nil
}

func (s *ServiceImpl) GetAllProjects(ctx context.Context) ([]Project, error) {
	projects, err := s.repo.GetAllProjects(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get projects")
	}

	return projects, nil
}

func (s *ServiceImpl) GetProjectByID(ctx context.Context, id int) (*Project, error) {
	project, err := s.repo.GetProjectByID(ctx, id)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get project")
	}

	return project, nil
}

func (s *ServiceImpl) GetProjectVersions(ctx context.Context, projectID int) ([]Version, error) {
	versions, err := s.repo.GetProjectVersions(ctx, projectID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get project versions")
	}

	return versions, nil
}

func (s *ServiceImpl) DeleteProjectVersionByIDAndProjectID(ctx context.Context, ID, projectID uint) error {
	version, err := s.repo.GetVersionByIDAndProjectID(ctx, ID, projectID)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return ErrNotFound
		}

		return errors.Wrapf(err, "Failed to retrieve project version with ID %d and project ID %d", ID, projectID)
	}

	return s.deleteProjectVersion(ctx, version)
}

func (s *ServiceImpl) deleteProjectVersion(ctx context.Context, version *Version) error {
	deleteVersionByIDTransaction, err := s.repo.DeleteVersionByIDTransaction(ctx, version.ID)
	if err != nil {
		return errors.Wrap(err, "Failed to begin delete project version transaction")
	}

	defer func() {
		if r := recover(); r != nil {
			deleteVersionByIDTransaction.Rollback()
		}
	}()

	if err := s.storage.DeleteObject(ctx, version.StorageKey); err != nil {
		deleteVersionByIDTransaction.Rollback()
		return errors.Wrap(err, "Failed to delete project version in S3")
	}

	if err := deleteVersionByIDTransaction.Commit().Error; err != nil {
		deleteVersionByIDTransaction.Rollback()
		return errors.Wrap(err, "Failed to commit delete project version transaction")
	}

	return nil
}

func (s *ServiceImpl) CreateLanguagesVersion(ctx context.Context, projectID int, name string) error {
	languagesResponse, err := s.Client.ListProjectLanguages(ctx, poedit.ListProjectLanguagesRequest{
		ID: projectID,
	})
	if err != nil {
		return errors.Wrap(err, "Failed to list project languages")
	}

	for _, language := range languagesResponse.Result.Languages {
		resp, err := s.Client.ExportProject(ctx, poedit.ExportProjectRequest{
			ID:       projectID,
			Language: language.Code,
			Type:     "key_value_json",
			Filters:  []string{"translated"},
		})
		if err != nil {
			return errors.Wrap(err, "Failed to export project")
		}

		d, err := http.Get(resp.Result.URL)
		if err != nil {
			return errors.Wrap(err, "Failed to download project language file")
		}
		defer d.Body.Close()

		if d.StatusCode != http.StatusOK {
			return errors.Errorf("Response code '%d' from download GET", d.StatusCode)
		}

		//upload to s3
		data, err := io.ReadAll(d.Body)
		if err != nil {
			return errors.Wrap(err, "Failed to read project language file")
		}

		reader := bytes.NewReader(data)
		key := fmt.Sprintf("%d/%s/%s.json", projectID, name, language.Code)
		err = s.storage.PutObject(ctx, key, reader, "application/json")
		if err != nil {
			return errors.Wrap(err, "Failed to upload project language file to S3")
		}
	}

	return nil
}
