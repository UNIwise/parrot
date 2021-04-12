package cache

import (
	"context"
	"errors"
	"time"
)

var (
	ErrCacheMiss = errors.New("Cache miss")
)

type Item struct {
	CreatedAt time.Time
	Checksum  string
	Data      interface{}
}

type Cache interface {
	GetProjectMeta(ctx context.Context, projectID int) (item *Item, err error)
	SetProjectMeta(ctx context.Context, projectID int, data interface{}) (checksum string, err error)
	ClearProjectMeta(ctx context.Context, projectID int) (err error)

	GetLanguage(ctx context.Context, projectID int, languageCode, format string) (item *Item, err error)
	SetLanguage(ctx context.Context, projectID int, languageCode, format string, data []byte) (checksum string, err error)
	ClearLanguage(ctx context.Context, projectID int, languageCode string) (err error)
	ClearProjectLanguages(ctx context.Context, projectID int) (err error)

	GetTTL() time.Duration

	PingContext(ctx context.Context) error
}
