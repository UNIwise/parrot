package cache

import (
	"context"
	"fmt"
	"time"

	redisCache "github.com/go-redis/cache/v8"
	"github.com/go-redis/redis/v8"

	"github.com/pkg/errors"
)

type RedisCache struct {
	rc  *redisCache.Cache
	ttl time.Duration
}

func NewRedisCache(rc *redis.Client, ttl time.Duration) *RedisCache {
	return &RedisCache{
		rc: redisCache.New(&redisCache.Options{
			Redis: rc,
		}),
		ttl: ttl,
	}
}

func (r *RedisCache) GetTranslation(ctx context.Context, projectID int, languageCode, format string) ([]byte, error) {
	key := r.key(projectID, languageCode, format)

	var data []byte
	err := r.rc.Get(ctx, key, data)
	if err == redis.Nil {
		return nil, ErrCacheMiss
	}
	if err != nil {
		return nil, errors.Wrapf(err, "Could not get cache data for key %s", key)
	}

	return data, nil
}

func (r *RedisCache) SetTranslation(ctx context.Context, projectID int, languageCode, format string, data []byte) error {
	key := r.key(projectID, languageCode, format)

	if err := r.rc.Set(&redisCache.Item{
		Ctx:            ctx,
		Key:            key,
		TTL:            r.ttl,
		Value:          data,
		SkipLocalCache: true,
	}); err != nil {
		return errors.Wrapf(err, "Error while setting cache data for key %s", key)
	}

	return nil
}

func (r *RedisCache) key(projectID int, languageCode, format string) string {
	return fmt.Sprintf("%d:%s:%s", projectID, languageCode, format)
}
