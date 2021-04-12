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
	"golang.org/x/sys/unix"
)

type FilesystemCache struct {
	dir string
	ttl time.Duration
}

func NewFilesystemCache(cacheDir string, ttl time.Duration) (*FilesystemCache, error) {
	err := os.MkdirAll(cacheDir, 0777)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to create cache directory")
	}

	return &FilesystemCache{
		dir: cacheDir,
		ttl: ttl,
	}, nil
}

func (r *FilesystemCache) SetProjectMeta(ctx context.Context, projectID int, meta interface{}) (item *Item, err error) {
	return nil, errors.New("Not implemented")
}

func (r *FilesystemCache) GetProjectMeta(ctx context.Context, projectID int) (checksum string, err error) {
	return "", errors.New("Not implemented")

}

func (r *FilesystemCache) ClearProjectMeta(ctx context.Context, projectID int) (err error) {
	return errors.New("Not implemented")
}

func (f *FilesystemCache) GetLanguage(ctx context.Context, projectID int, languageCode, format string) (*Item, error) {
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
	if err != nil {
		return nil, err
	}

	md5, err := ioutil.ReadFile(fmt.Sprintf("%s.md5", filePath))
	if err != nil {
		return nil, ErrCacheMiss
	}

	return &Item{
		Checksum:  string(md5),
		Data:      b,
		CreatedAt: s.ModTime(),
	}, nil
}

func (f *FilesystemCache) SetLanguage(ctx context.Context, projectID int, languageCode, format string, data []byte) (string, error) {
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

func (f *FilesystemCache) ClearLanguage(ctx context.Context, projectID int, languageCode string) error {
	prefix := fmt.Sprintf("%d_%s", projectID, languageCode)

	err := f.removeFilesWithPrefix(prefix)
	if err != nil {
		return errors.Wrapf(err, "Failed to remove cached translation files '%s/%s*'", f.dir, prefix)
	}

	return nil
}

func (f *FilesystemCache) ClearProjectLanguages(ctx context.Context, projectID int) error {
	prefix := fmt.Sprintf("%d_", projectID)

	err := f.removeFilesWithPrefix(prefix)
	if err != nil {
		return errors.Wrapf(err, "Failed to remove cached translation files '%s/%s*'", f.dir, prefix)
	}

	return nil
}

func (f *FilesystemCache) GetTTL() time.Duration {
	return f.ttl
}

func (f *FilesystemCache) PingContext(ctx context.Context) error {
	return unix.Access(f.dir, unix.W_OK)
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
