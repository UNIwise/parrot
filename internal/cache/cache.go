package cache

import "errors"

var (
	ErrCacheMiss = errors.New("Cache miss")
)

type Cache interface {
	GetTranslation(projectID int, languageCode, format string) ([]byte, error)
	SetTranslation(projectID int, languageCode, format string, contents []byte) error
}

type CacheImpl struct{}

func New() *CacheImpl {
	return &CacheImpl{}
}
