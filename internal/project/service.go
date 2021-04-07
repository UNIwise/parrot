package project

import (
	"context"
	"time"

	"io/ioutil"
	"net/http"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/uniwise/parrot/internal/cache"
	"github.com/uniwise/parrot/internal/poedit"
	"golang.org/x/sync/semaphore"
)

type Project struct {
	TTL      time.Duration
	Checksum string
	Meta     ProjectMeta
}

type ProjectMeta struct {
	Languages []ProjectMetaLanguage
}

type ProjectMetaLanguage struct {
	Code       string
	Updated    time.Time
	Percentage float64
}

type Language struct {
	TTL      time.Duration
	Checksum string
	Data     []byte
}

type Service interface {
	// Get meta data about a project.
	GetProjectMeta(ctx context.Context, projectID int) (meta *Project, err error)

	// Clear the cached meta data for a project.
	ClearProjectMeta(ctx context.Context, projectID int) (err error)

	// Get a language in a specific format from a project.
	GetLanguage(ctx context.Context, projectID int, languageCode, format string) (trans *Language, err error)

	// Clear the cached language for a project.
	ClearLanguage(ctx context.Context, projectID int, languageCode string) (err error)

	// Clear all cached languages for a project.
	ClearProjectLanguage(ctx context.Context, projectID int) (err error)
}

type ServiceImpl struct {
	Logger            *logrus.Entry
	Client            poedit.Client
	Cache             cache.Cache
	RenewalThreshold  time.Duration
	PreFetchSemaphore *semaphore.Weighted
}

func NewService(cli poedit.Client, cache cache.Cache, renewalThreshold time.Duration, entry *logrus.Entry) *ServiceImpl {
	return &ServiceImpl{
		Logger:            entry,
		Client:            cli,
		Cache:             cache,
		RenewalThreshold:  renewalThreshold,
		PreFetchSemaphore: semaphore.NewWeighted(1),
	}
}

func (s *ServiceImpl) GetProjectMeta(ctx context.Context, projectID int) (project *Project, err error) {
	res, err := s.Client.ListProjectLanguages(ctx, poedit.ListProjectLanguagesRequest{
		ID: projectID,
	})
	if err != nil {
		return nil, err
	}

	languages := make([]ProjectMetaLanguage, len(res.Result.Languages))
	for i, l := range res.Result.Languages {
		languages[i] = ProjectMetaLanguage{
			Code:       l.Code,
			Updated:    time.Now(),
			Percentage: float64(l.Percentage),
		}
	}

	return &Project{
		TTL:      s.Cache.GetTTL(),
		Checksum: "",
		Meta: ProjectMeta{
			Languages: languages,
		},
	}, nil
}

func (s *ServiceImpl) ClearProjectMeta(ctx context.Context, projectID int) error {
	return errors.New("Not implemented")
}

func (s *ServiceImpl) GetLanguage(ctx context.Context, projectID int, languageCode, format string) (*Language, error) {
	item, err := s.Cache.GetLanguage(ctx, projectID, languageCode, format)
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

				_, _, err := s.fetchAndCacheLanguage(context.Background(), projectID, languageCode, format)
				if err != nil {
					s.Logger.Errorf("Failed to pre-fetch language %s format %s for project %d", languageCode, format, projectID)
				}
			}()
		}

		return &Language{
			TTL:      s.Cache.GetTTL(),
			Checksum: item.Checksum,
			Data:     item.Data,
		}, nil
	}

	data, checksum, err := s.fetchAndCacheLanguage(ctx, projectID, languageCode, format)
	if err != nil {
		return nil, err
	}

	return &Language{
		TTL:      s.Cache.GetTTL(),
		Checksum: checksum,
		Data:     data,
	}, nil
}

func (s *ServiceImpl) ClearLanguage(ctx context.Context, projectID int, languageCode string) error {
	return errors.New("Not implemented")
}

func (s *ServiceImpl) ClearProjectLanguage(ctx context.Context, projectID int) error {
	return errors.New("Not implemented")
}

func (s *ServiceImpl) fetchAndCacheProjectMeta(ctx context.Context, projectID int) (*Project, string, error) {
	res, err := s.Client.ListProjectLanguages(ctx, poedit.ListProjectLanguagesRequest{
		ID: projectID,
	})
	if err != nil {
		return nil, "", err
	}

	checksum, err := s.Cache.SetProjectMeta(ctx, projectID, res)
	if err != nil {
		return nil, "", err
	}

	return res, checksum, nil
}

func (s *ServiceImpl) fetchAndCacheLanguage(ctx context.Context, projectID int, languageCode, format string) ([]byte, string, error) {
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

	checksum, err := s.Cache.SetLanguage(ctx, projectID, languageCode, format, data)
	if err != nil {
		return nil, "", err
	}

	return data, checksum, nil
}
