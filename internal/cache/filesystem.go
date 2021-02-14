package cache

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"time"
)

type FileCache struct {
	dir string
	ttl time.Duration
}

func NewFileCache(ttl time.Duration) (*FileCache, error) {
	dir, err := os.UserCacheDir()
	if err != nil {
		return nil, err
	}
	cacheDir := path.Join(dir, "parrot")
	err = os.MkdirAll(cacheDir, os.ModeDir)
	if err != nil {
		return nil, err
	}
	return &FileCache{
		dir: cacheDir,
		ttl: ttl,
	}, nil
}

func (f *FileCache) GetTranslation(projectID int, languageCode, format string) ([]byte, error) {
	s, err := os.Stat(f.key(projectID, languageCode, format))
	if os.IsNotExist(err) {
		return nil, ErrCacheMiss
	}
	if time.Since(s.ModTime()) > f.ttl {
		return nil, ErrCacheMiss
	}
	b, err := ioutil.ReadFile(f.key(projectID, languageCode, format))
	if os.IsNotExist(err) {
		return nil, ErrCacheMiss
	}
	return b, nil
}

func (f *FileCache) SetTranslation(projectID int, languageCode, format string, contents []byte) error {
	return ioutil.WriteFile(f.key(projectID, languageCode, format), contents, os.ModePerm)
}

func (f *FileCache) key(projectID int, languageCode, format string) string {
	return path.Join(f.dir, fmt.Sprintf("%d-%s-%s", projectID, languageCode, format))
}
