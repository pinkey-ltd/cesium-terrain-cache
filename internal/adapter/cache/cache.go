package cache

import (
	"container/list"
	"fmt"
	"log"
	"sync"
	"time"
)

/*
Item represents a single cache entry in the LRU (Least Recently Used) cache.
It contains the value of the cached item, the expiration time (if any),
and a reference to its corresponding list element in the LRU doubly linked list.
The reference to the list element allows the cache to efficiently update the position of the item
when it is accessed, ensuring the correct eviction of the least recently used items.
*/
type Item struct {
	value     string        // The value stored in the cache item
	expiresAt time.Time     // The expiration time for the cache item (if applicable)
	elem      *list.Element // A reference to the element in the LRU doubly linked list
}

/*
Cache represents an LRU (Least Recently Used) cache that stores key-value pairs.
The cache supports setting and getting items, with optional TTL (Time To Live) for each cache entry.
When the cache exceeds its maximum size, it evicts the least recently used item to make room for new entries.
Additionally, the cache can persist commands to an AOF (Append-Only File) for persistence between restarts.
*/
type Cache struct {
	mu      sync.RWMutex     // Mutex to ensure thread safety for concurrent access to the cache
	items   map[string]*Item // Map that stores cache items by their key
	lru     *list.List       // Doubly linked list used for tracking the access order of cache items (for LRU eviction)
	maxSize int              // Maximum number of items the cache can hold before eviction occurs
	persist *Persistence     // Optional persistence mechanism for appending commands to a file
}

/*
NewCache initializes and returns a new Cache instance with the specified max size.
The cache starts with an empty linked list and an empty item map, and optionally supports persistence.
*/
func NewCache(maxSize int, persist *Persistence) *Cache {
	return &Cache{
		items:   make(map[string]*Item), // Initialize the map to store cache items
		lru:     list.New(),             // Initialize the doubly linked list to track the LRU order
		maxSize: maxSize,                // Set the maximum size for the cache
		persist: persist,                // Optionally set the persistence mechanism
	}
}

/*
Set stores a key-value pair in the cache, optionally with a TTL (Time To Live).
If the cache exceeds the maximum size, the least recently used (LRU) item is evicted.
If persistence is enabled, the SET command is appended to the AOF file.
*/
func (c *Cache) Set(key string, value string, ttl time.Duration, replaying bool) {
	c.mu.Lock()         // Lock the cache to ensure thread-safe access
	defer c.mu.Unlock() // Unlock once the operation is complete

	var expiresAt time.Time
	if ttl != 0 {
		expiresAt = time.Now().Add(ttl) // Set the expiration time if TTL is provided
	}

	// Check if the key already exists in the cache
	if item, exists := c.items[key]; exists {
		// Update the value and expiration time of the existing item
		item.value = value
		item.expiresAt = expiresAt
		// Move the item to the front of the LRU list to mark it as recently used
		c.lru.MoveToFront(item.elem)
	} else {
		// Create a new item and insert it at the front of the LRU list
		elem := c.lru.PushFront(key)
		c.items[key] = &Item{
			value:     value,
			expiresAt: expiresAt,
			elem:      elem, // Store the list element reference with the item
		}
	}

	// Evict the least recently used item if the cache exceeds the max size
	if c.lru.Len() > c.maxSize {
		oldest := c.lru.Back() // Get the least recently used item (the oldest in the list)
		if oldest != nil {
			delete(c.items, oldest.Value.(string)) // Remove the item from the map
			c.lru.Remove(oldest)                   // Remove the item from the LRU list
		}
	}

	// Persist the command in the AOF file if persistence is enabled and not replaying
	if c.persist != nil && !replaying {
		cmd := fmt.Sprintf("SET %s %s", key, value)
		err := c.persist.Append(cmd)
		if err != nil {
			log.Println("Failed to append to AOF file:", err)
		}
	}
}

/*
Get retrieves the value associated with a key from the cache.
If the key exists and has not expired, it is moved to the front of the LRU list.
If the key does not exist or has expired, it is removed from the cache and not returned.
*/
func (c *Cache) Get(key string) (string, bool) {
	c.mu.Lock()         // Lock the cache to ensure thread-safe access
	defer c.mu.Unlock() // Unlock once the operation is complete

	// Check if the key exists in the cache
	item, exists := c.items[key]
	if !exists {
		return "", false // Return false if the key doesn't exist
	}

	// Check if the item has expired
	if !item.expiresAt.IsZero() && time.Now().After(item.expiresAt) {
		// If expired, remove it from the LRU list and delete it from the cache
		c.lru.Remove(item.elem)
		delete(c.items, key)
		return "", false
	}

	// Move the item to the front of the LRU list to mark it as recently used
	c.lru.MoveToFront(item.elem)

	// Return the value of the item along with a true flag indicating success
	return item.value, true
}

/*
StartCleaningServer periodically runs a cleanup task that removes expired items from the cache.
This task is triggered every minute and ensures the cache does not hold stale data.
*/
func (c *Cache) StartCleaningServer() {
	ticker := time.NewTicker(time.Minute * 1) // Set the cleanup interval to 1 minute
	defer ticker.Stop()

	// Run the cleanup process at regular intervals
	for range ticker.C {
		log.Println("Running cache cleanup")
		c.cleanExpiredItems() // Perform the cleanup of expired items
	}
}

/*
cleanExpiredItems iterates through the cache and removes any expired items.
It checks each item in the LRU list, and if the item has expired, it is removed from both the list and the cache.
Additionally, if the cache exceeds its maximum size, the least recently used item is evicted.
*/
func (c *Cache) cleanExpiredItems() {
	c.mu.Lock()         // Lock the cache for thread-safe access
	defer c.mu.Unlock() // Unlock once the operation is complete

	// Iterate through the LRU list from the back (oldest) to the front (most recently used)
	for e := c.lru.Back(); e != nil; e = e.Prev() {
		item := c.items[e.Value.(string)]

		// Skip items that do not have an expiration time
		if item.expiresAt.IsZero() {
			continue
		}

		// Remove expired items
		if time.Now().After(item.expiresAt) {
			c.lru.Remove(e)                   // Remove the expired item from the LRU list
			delete(c.items, e.Value.(string)) // Remove the expired item from the cache map
			log.Printf("Removed expired item: %s\n", e.Value)
		}
	}

	// If the cache exceeds the maximum size, evict the least recently used item
	if c.lru.Len() > c.maxSize {
		oldest := c.lru.Back()
		if oldest != nil {
			delete(c.items, oldest.Value.(string))
			c.lru.Remove(oldest)
		}
	}
}
