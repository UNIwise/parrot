package cache

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"time"

	"github.com/pkg/errors"
)

const (
	cacheSubdir = "parrot"
)

type FilesystemCache struct {
	dir string
	ttl time.Duration
}

func NewFilesystemCache(ttl time.Duration) (*FilesystemCache, error) {
	dir, err := os.UserCacheDir()
	if err != nil {
		return nil, errors.Wrap(err, "Failed to determin OS user cache directory")
	}

	cacheDir := path.Join(dir, cacheSubdir)

	err = os.MkdirAll(cacheDir, os.ModeDir)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to create cacge directory")
	}

	return &FilesystemCache{
		dir: cacheDir,
		ttl: ttl,
	}, nil
}

func (f *FilesystemCache) GetTranslation(ctx context.Context, projectID int, languageCode, format string) ([]byte, error) {
	filename := f.key(projectID, languageCode, format)

	s, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return nil, ErrCacheMiss
	} else if err != nil {
		return nil, errors.Wrap(err, "Failed to get cached file state from OS")
	}

	if time.Since(s.ModTime()) > f.ttl {
		if err := os.Remove(filename); err != nil {
			return nil, errors.Wrap(err, "Failed to remove expired cache file")
		}

		return nil, ErrCacheMiss
	}

	b, err := ioutil.ReadFile(filename)
	if os.IsNotExist(err) {
		return nil, ErrCacheMiss
	}

	return b, nil
}

func (f *FilesystemCache) SetTranslation(ctx context.Context, projectID int, languageCode, format string, data []byte) error {
	return ioutil.WriteFile(
		f.key(
			projectID,
			languageCode,
			format,
		),
		data,
		os.ModePerm,
	)
}

func (f *FilesystemCache) key(projectID int, languageCode, format string) string {
	return path.Join(
		f.dir,
		fmt.Sprintf(
			"%d-%s-%s",
			projectID,
			languageCode,
			format,
		),
	)
}
