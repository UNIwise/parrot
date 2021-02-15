package cache

import (
	"context"
	"errors"
)

var (
	ErrCacheMiss = errors.New("Cache miss")
)

type Cache interface {
	GetTranslation(ctx context.Context, projectID int, languageCode, format string) (data []byte, err error)
	SetTranslation(ctx context.Context, projectID int, languageCode, format string, data []byte) (err error)
	PurgeTranslation(ctx context.Context, projectID int, languageCode string) (err error)
	PurgeProject(ctx context.Context, projectID int) (err error)
}
