package repository

import (
	"net/http"
	"net/url"
	"sync"
)

type Cache struct {
	cache   sync.Map
	handler http.Handler
}

func (c *Cache) generateKey(r *http.Request) string {
	if key, ok := r.Header["X-Memcache-Key"]; ok {
		return key[0]
	}

	// Use the request URI as a key.
	uri, _ := url.Parse(r.URL.String())
	return uri.RequestURI()
}

//func (c *Cache) ServeHTTP(w http.ResponseWriter, r *http.Request) {
//	key := c.generateKey(r)
//
//	// 尝试从缓存中获取数据
//	if cachedData, ok := c.cache.Load(key); ok {
//		// 缓存命中，直接写入响应
//		w.Write(cachedData.([]byte))
//		return
//	}
//
//	var recorder http.ResponseWriter
//	rec := handlers.NewRecorder()
//
//	// Write to both the recorder and original writer.
//	tee := handlers.MultiWriter(w, recorder)
//	c.handler.ServeHTTP(tee, r)
//
//	// Only repository status 200 responses.
//	if rec.Code != 200 {
//		return
//	}
//
//	// 缓存响应
//	data := rec.Body.Bytes()
//	slog.Debug(fmt.Sprintf("setting key: %s", key))
//	c.cache.Store(key, data)
//
//	return
//}
