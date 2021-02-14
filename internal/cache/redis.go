package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
)

type RedisCache struct {
	rc  *redis.Client
	ttl time.Duration
}

func NewRedisCache(rc *redis.Client, ttl time.Duration) *RedisCache {
	return &RedisCache{
		rc:  rc,
		ttl: ttl,
	}
}

func (r *RedisCache) GetTranslation(projectID int, languageCode, format string) ([]byte, error) {
	ctx := context.TODO()
	b, err := r.rc.Get(ctx, r.key(projectID, languageCode, format)).Bytes()
	if err == redis.Nil {
		return nil, ErrCacheMiss
	} else if err != nil {
		return nil, errors.Wrap(err, "Could not contents of cache")
	}
	return b, nil
}

func (r *RedisCache) SetTranslation(projectID int, languageCode, format string, contents []byte) error {
	ctx := context.TODO()
	err := r.rc.Set(ctx, r.key(projectID, languageCode, format), contents, r.ttl).Err()
	if err != nil {
		return errors.Wrap(err, "Error while setting cache value")
	}
	return nil
}

func (r *RedisCache) key(projectID int, languageCode, format string) string {
	return fmt.Sprintf("%d-%s-%s", projectID, languageCode, format)
}
