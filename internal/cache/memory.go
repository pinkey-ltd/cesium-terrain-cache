package cache

import (
	"fmt"
	"sync"
	"time"
)

type CacheItem struct {
	value      interface{}
	expiration int64
}

type InMemoryCache struct {
	items sync.Map
}

func NewInMemoryCache() *InMemoryCache {
	return &InMemoryCache{}
}

func (c *InMemoryCache) Set(key string, value interface{}, duration time.Duration) {
	expiration := time.Now().Add(duration).UnixNano()
	c.items.Store(key, CacheItem{
		value:      value,
		expiration: expiration,
	})
}

func (c *InMemoryCache) Get(key string) (interface{}, bool) {
	item, found := c.items.Load(key)
	if !found {
		return nil, false
	}

	cacheItem := item.(CacheItem)
	if time.Now().UnixNano() > cacheItem.expiration {
		c.items.Delete(key)
		return nil, false
	}

	return cacheItem.value, true
}

func (c *InMemoryCache) Delete(key string) {
	c.items.Delete(key)
}

func (c *InMemoryCache) Cleanup() {
	c.items.Range(func(key, value interface{}) bool {
		cacheItem := value.(CacheItem)
		if time.Now().UnixNano() > cacheItem.expiration {
			c.items.Delete(key)
		}
		return true
	})
}

func main() {
	cache := NewInMemoryCache()

	// 设置缓存项，过期时间为2秒
	cache.Set("foo", "bar", 2*time.Second)

	// 获取缓存项
	if value, found := cache.Get("foo"); found {
		fmt.Println("Found:", value)
	} else {
		fmt.Println("Not found")
	}

	// 等待3秒后再次获取
	time.Sleep(3 * time.Second)
	if value, found := cache.Get("foo"); found {
		fmt.Println("Found:", value)
	} else {
		fmt.Println("Not found")
	}

	// 清理过期缓存
	cache.Cleanup()
}
