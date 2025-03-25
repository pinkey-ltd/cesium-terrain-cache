package cache

import (
	"fmt"
	"github.com/pinkey-ltd/cesium-terrain-server/handlers"
	"log/slog"
	"net/http"
	"net/url"
	"sync"
)

type Cache struct {
	cache   sync.Map
	handler http.Handler
	Limit   handlers.Bytes
	limiter handlers.LimiterFactory
}

// NewCache creates a new HTTP caching handler with a storage limit and a response limiter factory.
func NewCache(handler http.Handler, limit handlers.Bytes, limiter handlers.LimiterFactory) http.Handler {
	return &Cache{
		cache:   sync.Map{},
		handler: handler,
		Limit:   limit,
		limiter: limiter,
	}
}

func (c *Cache) generateKey(r *http.Request) string {
	if key, ok := r.Header["X-Memcache-Key"]; ok {
		return key[0]
	}

	// Use the request URI as a key.
	uri, _ := url.Parse(r.URL.String())
	return uri.RequestURI()
}

func (c *Cache) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	key := c.generateKey(r)

	// 尝试从缓存中获取数据
	if cachedData, ok := c.cache.Load(key); ok {
		// 缓存命中，直接写入响应
		w.Write(cachedData.([]byte))
		return
	}

	var limiter handlers.ResponseLimiter
	var recorder http.ResponseWriter
	rec := handlers.NewRecorder()

	// If a limiter is provided, wrap the recorder with it.
	if c.limiter != nil {
		limiter = c.limiter(rec, c.Limit)
		recorder = limiter
	} else {
		recorder = rec
	}

	// Write to both the recorder and original writer.
	tee := handlers.MultiWriter(w, recorder)
	c.handler.ServeHTTP(tee, r)

	// Only cache status 200 responses.
	if rec.Code != 200 {
		return
	}

	// If the cache limit has been exceeded, don't proceed to cache the
	// response.
	if limiter != nil && limiter.LimitExceeded() {
		slog.Debug(fmt.Sprintf("cache limit exceeded for %s", r.URL.String()))
		return
	}

	// 缓存响应
	data := rec.Body.Bytes()
	slog.Debug(fmt.Sprintf("setting key: %s", key))
	c.cache.Store(key, data)

	return
}
