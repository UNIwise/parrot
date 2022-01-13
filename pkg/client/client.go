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

// Client is a client for fetching translations from parrot.
type Client interface {
	// GetTranslation returns a map of terms for the given language.
	GetTranslation(ctx context.Context, language string) (map[string]string, error)
}

// CachedClientImpl is a client that caches translations in memory.
type CachedClientImpl struct {
	httpClient *resty.Client
	projectID  int
	cache      map[string]cacheItem
	mutex      *sync.Mutex
	ttl        time.Duration
}

// NewCachedClient creates a new client that caches translations in memory.
func NewCachedClient(endpoint string, project int, ttl time.Duration) *CachedClientImpl {
	return &CachedClientImpl{
		httpClient: resty.New().SetHostURL(endpoint),
		projectID:  project,
		cache:      map[string]cacheItem{},
		mutex:      &sync.Mutex{},
		ttl:        ttl,
	}
}

// cacheItem is a single cached translation, as it appears in the cache.
type cacheItem struct {
	Etag    string
	Data    map[string]string
	Expires time.Time
}

// response is the response type returned from parrot, with the json format.
type response []struct {
	Term       string `json:"term"`
	Definition string `json:"definition"`
}

// GetTranslation returns a map of terms for the given language. It will return cached data
// if it has not expired, otherwise it will fetch the data from the upstream.
func (c *CachedClientImpl) GetTranslation(ctx context.Context, language string) (map[string]string, error) {
	key := cacheKey(c.projectID, language)
	// Lock here to avoid multiple concurrent requests to upstream
	c.mutex.Lock()
	defer c.mutex.Unlock()
	cached, exists := c.cache[key]
	if exists && cached.Expires.After(time.Now()) {
		return cached.Data, nil
	}

	// Fetch from upstream and store in the cache
	request := c.httpClient.R().
		SetContext(ctx).
		SetPathParams(map[string]string{
			"project":  strconv.Itoa(c.projectID),
			"language": language,
		}).
		SetQueryParam("format", "json").
		SetResult(response{})

	if exists {
		request = request.SetHeader("If-None-Match", cached.Etag)
	}

	resp, err := request.Get("/v1/project/{project}/language/{language}")
	if err != nil {
		return nil, errors.Wrap(err, "Failed to fetch translation")
	}

	switch resp.StatusCode() {
	case http.StatusNotModified:
		// Handle 304 Not Modified by renewing the cached item
		c.cache[key] = cacheItem{
			Etag:    cached.Etag,
			Data:    cached.Data,
			Expires: time.Now().Add(c.ttl),
		}

		return cached.Data, nil
	case http.StatusOK:
		// Handle 200 OK by storing the response in the cache
		data, ok := resp.Result().(*response)
		if !ok || data == nil {
			return nil, errors.Errorf("Failed to parse response: '%s'", data)
		}
		result := mapResponse(*data)
		c.cache[key] = cacheItem{
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

// mapResponse converts the response from parrot into a map of terms.
func mapResponse(resp response) map[string]string {
	result := map[string]string{}

	for _, item := range resp {
		result[item.Term] = item.Definition
	}

	return result
}
