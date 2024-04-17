package in_memory_cache

import (
	"sync"
	"time"
)

// CacheEntry представляет запись в кеше
type CacheEntry struct {
	Data      interface{} // Результат запроса
	ExpiresAt time.Time   // Время истечения срока действия записи в кеше
}

// InMemoryCache представляет кеш в памяти
type InMemoryCache struct {
	cache map[string]CacheEntry // Карта для хранения записей в кеше
	mu    sync.RWMutex          // Мьютекс для защиты доступа к кешу
}

// NewInMemoryCache создает новый экземпляр кеша в памяти
func NewInMemoryCache() *InMemoryCache {
	return &InMemoryCache{
		cache: make(map[string]CacheEntry),
	}
}

// Get возвращает данные из кеша для указанного ключа
func (c *InMemoryCache) Get(key string) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	entry, ok := c.cache[key]
	if !ok || entry.ExpiresAt.Before(time.Now()) {
		// Запись отсутствует в кеше или истек срок ее действия
		return nil, false
	}
	return entry.Data, true
}

// Set устанавливает данные в кеш для указанного ключа
func (c *InMemoryCache) Set(key string, data interface{}, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.cache[key] = CacheEntry{
		Data:      data,
		ExpiresAt: time.Now().Add(ttl),
	}
}

// Clear очищает кеш
func (c *InMemoryCache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.cache = make(map[string]CacheEntry)
}

// CleanupExpiredEntries сканирует кеш и удаляет истекшие записи
func (c *InMemoryCache) CleanupExpiredEntries() {
	c.mu.Lock()
	defer c.mu.Unlock()

	for key, entry := range c.cache {
		if entry.ExpiresAt.Before(time.Now()) {
			delete(c.cache, key)
		}
	}
}
