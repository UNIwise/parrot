package cache

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/pkg/errors"
)

type FilesystemCache struct {
	dir string
	ttl time.Duration
}

func NewFilesystemCache(cacheDir string, ttl time.Duration) (*FilesystemCache, error) {
	err := os.MkdirAll(cacheDir, 0o777)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to create cache directory")
	}

	return &FilesystemCache{
		dir: cacheDir,
		ttl: ttl,
	}, nil
}

func (f *FilesystemCache) GetTranslation(ctx context.Context, projectID int, languageCode, format string) (*CacheItem, error) {
	filePath := f.filePath(projectID, languageCode, format)

	info, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		return nil, ErrCacheMiss
	}
	if err != nil {
		return nil, errors.Wrap(err, "Failed to get cached file state from OS")
	}

	if time.Since(info.ModTime()) > f.ttl {
		return nil, ErrCacheMiss
	}

	b, err := ioutil.ReadFile(filePath)
	if os.IsNotExist(err) {
		return nil, ErrCacheMiss
	}
	if err != nil {
		return nil, err
	}

	md5, err := ioutil.ReadFile(fmt.Sprintf("%s.md5", filePath))
	if err != nil {
		return nil, ErrCacheMiss
	}

	return &CacheItem{
		Checksum:  string(md5),
		Data:      b,
		CreatedAt: info.ModTime(),
	}, nil
}

func (f *FilesystemCache) SetTranslation(ctx context.Context, projectID int, languageCode, format string, data []byte) (string, error) {
	if err := ioutil.WriteFile(
		f.filePath(projectID, languageCode, format),
		data,
		os.ModePerm,
	); err != nil {
		return "", err
	}

	hashBytes := md5.Sum(data)
	hash := hex.EncodeToString(hashBytes[:])

	if err := ioutil.WriteFile(
		f.md5Path(projectID, languageCode, format),
		[]byte(hash),
		os.ModePerm,
	); err != nil {
		return "", err
	}

	return hash, nil
}

func (f *FilesystemCache) PurgeTranslation(ctx context.Context, projectID int, languageCode string) error {
	prefix := fmt.Sprintf("%d_%s", projectID, languageCode)

	err := f.removeFilesWithPrefix(prefix)
	if err != nil {
		return errors.Wrapf(err, "Failed to remove cached translation files '%s/%s*'", f.dir, prefix)
	}

	return nil
}

func (f *FilesystemCache) PurgeProject(ctx context.Context, projectID int) error {
	prefix := fmt.Sprintf("%d_", projectID)

	err := f.removeFilesWithPrefix(prefix)
	if err != nil {
		return errors.Wrapf(err, "Failed to remove cached translation files '%s/%s*'", f.dir, prefix)
	}

	return nil
}

func (f *FilesystemCache) filePath(projectID int, languageCode, format string) string {
	return path.Join(
		f.dir,
		f.filename(projectID, languageCode, format),
	)
}

func (f *FilesystemCache) md5Path(projectID int, languageCode, format string) string {
	return path.Join(
		f.dir,
		fmt.Sprintf(
			"%s.md5",
			f.filename(projectID, languageCode, format),
		),
	)
}

func (f *FilesystemCache) filename(projectID int, languageCode, format string) string {
	return fmt.Sprintf("%d_%s_%s", projectID, languageCode, format)
}

func (f *FilesystemCache) removeFilesWithPrefix(prefix string) error {
	return filepath.Walk(f.dir, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			if strings.HasPrefix(info.Name(), prefix) {
				if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
					return err
				}
			}
		}

		return nil
	})
}

func (f *FilesystemCache) GetTTL() time.Duration {
	return f.ttl
}
