package cache

import (
	"context"
	"errors"
	"time"
)

var (
	ErrCacheMiss = errors.New("Cache miss")
)

type CacheItem struct {
	CreatedAt time.Time
	Checksum  string
	Data      []byte
}

type Cache interface {
	SetProjectMeta(ctx context.Context, projectID int, meta interface{}) (item *CacheItem, err error)
	GetProjectMeta(ctx context.Context, projectID int) (checksum string, err error)
	ClearProjectMeta(ctx context.Context, projectID int) (err error)

	GetLanguage(ctx context.Context, projectID int, languageCode, format string) (item *CacheItem, err error)
	SetLanguage(ctx context.Context, projectID int, languageCode, format string, data []byte) (checksum string, err error)
	ClearLanguage(ctx context.Context, projectID int, languageCode string) (err error)
	ClearProjectLanguages(ctx context.Context, projectID int) (err error)

	GetTTL() time.Duration
	PingContext(ctx context.Context) error
}
