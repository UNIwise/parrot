package cache

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"time"

	"github.com/pkg/errors"
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

func (f *FilesystemCache) GetTranslation(ctx context.Context, projectID int, languageCode, format string) ([]byte, string, error) {
	filePath := f.filePath(projectID, languageCode, format)

	s, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		return nil, "", ErrCacheMiss
	}
	if err != nil {
		return nil, "", errors.Wrap(err, "Failed to get cached file state from OS")
	}

	if time.Since(s.ModTime()) > f.ttl {
		return nil, "", ErrCacheMiss
	}

	b, err := ioutil.ReadFile(filePath)
	if os.IsNotExist(err) {
		return nil, "", ErrCacheMiss
	}
	if err != nil {
		return nil, "", err
	}

	md5, err := ioutil.ReadFile(fmt.Sprintf("%s.md5", filePath))
	if err != nil {
		return nil, "", ErrCacheMiss
	}

	return b, string(md5), nil
}

func (f *FilesystemCache) SetTranslation(ctx context.Context, projectID int, languageCode, format string, data []byte) (string, error) {
	filePath := f.filePath(projectID, languageCode, format)

	if err := ioutil.WriteFile(
		filePath,
		data,
		os.ModePerm,
	); err != nil {
		return "", err
	}

	hashBytes := md5.Sum(data)
	hash := hex.EncodeToString(hashBytes[:])

	if err := ioutil.WriteFile(
		fmt.Sprintf("%s.md5", filePath),
		[]byte(hash),
		os.ModePerm,
	); err != nil {
		return "", err
	}

	return hash, nil
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
