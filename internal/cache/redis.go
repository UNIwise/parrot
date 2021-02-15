package cache

import (
	"context"
	"crypto/md5"
	"encoding/hex"
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

type RedisCacheItem struct {
	Hash string
	Data []byte
}

func NewRedisCache(rc *redis.Client, ttl time.Duration) *RedisCache {
	return &RedisCache{
		rc: redisCache.New(&redisCache.Options{
			Redis: rc,
		}),
		ttl: ttl,
	}
}

func (r *RedisCache) GetTranslation(ctx context.Context, projectID int, languageCode, format string) ([]byte, string, error) {
	key := r.key(projectID, languageCode, format)

	var item RedisCacheItem
	err := r.rc.Get(ctx, key, item)
	if err == redis.Nil {
		return nil, "", ErrCacheMiss
	}
	if err != nil {
		return nil, "", errors.Wrapf(err, "Could not get cache data for key %s", key)
	}

	return item.Data, item.Hash, nil
}

func (r *RedisCache) SetTranslation(ctx context.Context, projectID int, languageCode, format string, data []byte) (string, error) {
	key := r.key(projectID, languageCode, format)

	hashBytes := md5.Sum(data)
	hash := hex.EncodeToString(hashBytes[:])

	if err := r.rc.Set(&redisCache.Item{
		Ctx: ctx,
		Key: key,
		TTL: r.ttl,
		Value: RedisCacheItem{
			Hash: hash,
			Data: data,
		},
		SkipLocalCache: true,
	}); err != nil {
		return "", errors.Wrapf(err, "Error while setting cache data for key %s", key)
	}

	return "", nil
}

func (f *RedisCache) PurgeTranslation(ctx context.Context, projectID int, languageCode string) (err error) {
	return errors.New("Not implemented")
}

func (f *RedisCache) PurgeProject(ctx context.Context, projectID int) (err error) {
	return errors.New("Not implemented")
}

func (r *RedisCache) key(projectID int, languageCode, format string) string {
	return fmt.Sprintf("%d:%s:%s", projectID, languageCode, format)
}
