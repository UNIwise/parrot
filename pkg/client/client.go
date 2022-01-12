package client

import (
	"context"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/pkg/errors"
	"gopkg.in/resty.v1"
)

type Client interface {
	GetTranslation(ctx context.Context, language string) (map[string]string, error)
}

type CachedClientImpl struct {
	httpClient *resty.Client
	projectID  int
	cache      map[string]CacheItem
	mutex      *sync.Mutex
	ttl        time.Duration
}

func NewCachedClient(endpoint string, project int, ttl time.Duration) *CachedClientImpl {
	return &CachedClientImpl{
		httpClient: resty.New().SetHostURL(endpoint),
		projectID:  project,
		cache:      map[string]CacheItem{},
		mutex:      &sync.Mutex{},
		ttl:        ttl,
	}
}

type CacheItem struct {
	Etag    string
	Data    map[string]string
	Expires time.Time
}

type Response []struct {
	Term       string `json:"term"`
	Definition string `json:"definition"`
}

func (c *CachedClientImpl) GetTranslation(ctx context.Context, language string) (map[string]string, error) {
	key := cacheKey(c.projectID, language)
	// Lock here to avoid multiple concurrent requests to upstream
	c.mutex.Lock()
	defer c.mutex.Unlock()
	cached, exists := c.cache[key]
	if exists && cached.Expires.After(time.Now()) {
		return cached.Data, nil
	}

	request := c.httpClient.R().
		SetContext(ctx).
		SetPathParams(map[string]string{
			"project":  strconv.Itoa(c.projectID),
			"language": language,
		}).
		SetQueryParam("format", "json").
		SetResult(Response{})

	if exists {
		request = request.SetHeader("If-None-Match", cached.Etag)
	}

	resp, err := request.Get("/v1/project/{project}/language/{language}")
	if err != nil {
		return nil, errors.Wrap(err, "Failed to get translation")
	}

	switch resp.StatusCode() {
	case http.StatusNotModified:
		// Renew cached item
		c.cache[key] = CacheItem{
			Etag:    cached.Etag,
			Data:    cached.Data,
			Expires: time.Now().Add(c.ttl),
		}

		return cached.Data, nil
	case http.StatusOK:
		// Cache response
		data, ok := resp.Result().(*Response)
		if !ok || data == nil {
			return nil, errors.Errorf("Failed to parse response: '%s'", data)
		}
		result := mapResponse(*data)
		c.cache[key] = CacheItem{
			Etag:    resp.Header().Get("Etag"),
			Data:    result,
			Expires: time.Now().Add(c.ttl),
		}

		return result, nil
	default:
		return nil, errors.Errorf("Unexpected response code: %d", resp.StatusCode())
	}
}

func cacheKey(projectID int, language string) string {
	return "translation:" + strconv.Itoa(projectID) + ":" + language
}

func mapResponse(resp Response) map[string]string {
	result := map[string]string{}

	for _, item := range resp {
		result[item.Term] = item.Definition
	}

	return result
}
