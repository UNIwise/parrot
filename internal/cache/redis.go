package cache

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	redisCache "github.com/go-redis/cache/v8"
	"github.com/go-redis/redis/v8"

	"github.com/pkg/errors"
)

const (
	redisScanCount = 10
)

type RedisCache struct {
	c   *redis.Client
	rc  *redisCache.Cache
	ttl time.Duration
}

type RedisCacheItem struct {
	Hash string
	Data []byte
}

func NewRedisCache(c *redis.Client, ttl time.Duration) *RedisCache {
	return &RedisCache{
		c: c,
		rc: redisCache.New(&redisCache.Options{
			Redis: c,
		}),
		ttl: ttl,
	}
}

func (r *RedisCache) GetTranslation(ctx context.Context, projectID int, languageCode, format string) ([]byte, string, error) {
	key := r.key(projectID, languageCode, format)

	var item RedisCacheItem
	err := r.rc.Get(ctx, key, &item)
	if err != nil {
		if strings.Contains(err.Error(), "key is missing") {
			return nil, "", ErrCacheMiss
		}

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

	return hash, nil
}

func (r *RedisCache) PurgeTranslation(ctx context.Context, projectID int, languageCode string) error {
	pattern := fmt.Sprintf("%d:%s:*", projectID, languageCode)

	if err := r.deleteKeysMatching(ctx, pattern); err != nil {
		return errors.Wrapf(err, "Failed to remove cached language '%s' for project '%d'", languageCode, projectID)
	}

	return nil
}

func (r *RedisCache) PurgeProject(ctx context.Context, projectID int) error {
	pattern := fmt.Sprintf("%d:*", projectID)

	if err := r.deleteKeysMatching(ctx, pattern); err != nil {
		return errors.Wrapf(err, "Failed to remove cached project '%d'", projectID)
	}

	return nil
}

func (r *RedisCache) key(projectID int, languageCode, format string) string {
	return fmt.Sprintf("%d:%s:%s", projectID, languageCode, format)
}

func (r *RedisCache) deleteKeysMatching(ctx context.Context, pattern string) error {
	keys, err := r.getKeysMatching(ctx, pattern)
	if err != nil {
		return err
	}

	if err := r.c.Del(ctx, keys...).Err(); err != nil {
		return errors.Wrapf(err, "Failed to remove redis keys matching '%s'", pattern)
	}

	return nil
}

func (r *RedisCache) getKeysMatching(ctx context.Context, pattern string) ([]string, error) {
	var allKeys []string

	var cursor uint64
	for {
		var keys []string
		var err error

		keys, cursor, err = r.c.Scan(
			ctx,
			cursor,
			pattern,
			redisScanCount,
		).Result()
		if err != nil {
			return nil, errors.Wrapf(err, "Failed to retrieve keys from redis matching pattern '%s'", pattern)
		}

		if cursor == 0 {
			break
		}

		allKeys = append(allKeys, keys...)
	}

	return allKeys, nil
}
