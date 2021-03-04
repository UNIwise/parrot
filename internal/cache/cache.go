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
	GetTranslation(ctx context.Context, projectID int, languageCode, format string) (item *CacheItem, err error)
	SetTranslation(ctx context.Context, projectID int, languageCode, format string, data []byte) (checksum string, err error)
	PurgeTranslation(ctx context.Context, projectID int, languageCode string) (err error)
	PurgeProject(ctx context.Context, projectID int) (err error)
	GetTTL() time.Duration
}
