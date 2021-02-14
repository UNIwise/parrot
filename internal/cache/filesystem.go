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
		return nil, errors.Wrap(err, "Failed to create cache directory")
	}

	return &FilesystemCache{
		dir: cacheDir,
		ttl: ttl,
	}, nil
}

func (f *FilesystemCache) GetTranslation(ctx context.Context, projectID int, languageCode, format string) ([]byte, error) {
	filePath := f.filePath(projectID, languageCode, format)

	s, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		return nil, ErrCacheMiss
	}
	if err != nil {
		return nil, errors.Wrap(err, "Failed to get cached file state from OS")
	}

	if time.Since(s.ModTime()) > f.ttl {
		return nil, ErrCacheMiss
	}

	b, err := ioutil.ReadFile(filePath)
	if os.IsNotExist(err) {
		return nil, ErrCacheMiss
	}

	return b, nil
}

func (f *FilesystemCache) SetTranslation(ctx context.Context, projectID int, languageCode, format string, data []byte) error {
	filePath := f.filePath(projectID, languageCode, format)

	return ioutil.WriteFile(
		filePath,
		data,
		os.ModePerm,
	)
}

func (f *FilesystemCache) PurgeTranslation(ctx context.Context, projectID int, languageCode string) (err error) {
	return errors.New("Not implemented")
}

func (f *FilesystemCache) PurgeProject(ctx context.Context, projectID int) (err error) {
	return errors.New("Not implemented")
}

func (f *FilesystemCache) filePath(projectID int, languageCode, format string) string {
	return path.Join(
		f.dir,
		f.filename(projectID, languageCode, format),
	)
}

func (f *FilesystemCache) filename(projectID int, languageCode, format string) string {
	return fmt.Sprintf("%d_%s_%s", projectID, languageCode, format)
}
