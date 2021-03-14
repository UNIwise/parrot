package project

import (
	"context"
	"time"

	"io/ioutil"
	"net/http"

	gosundheit "github.com/AppsFlyer/go-sundheit"
	"github.com/AppsFlyer/go-sundheit/checks"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/uniwise/parrot/internal/cache"
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
	c, err := checks.NewPingCheck("cache", s.Cache, time.Second*1)
	if err != nil {
		return errors.Wrap(err, "Failed to instantiate cache healthcheck")
	}

	if err := h.RegisterCheck(&gosundheit.Config{Check: c, ExecutionPeriod: time.Second * 10}); err != nil {
		return errors.Wrap(err, "Failed to register cache healthcheck")
	}

	return nil
}
