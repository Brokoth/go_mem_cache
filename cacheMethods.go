package gomemcache

import (
	"errors"
	"reflect"
	"time"
)

func (cache *Cache) Get(key interface{}) (value interface{}, err error) {
	var currentTime = time.Now().UTC()
	cacheEntry, ok := readMapEntry(cache, key)

	if !ok {
		return nil, errors.New("key not found")
	}

	if cacheEntry.Expires && (cacheEntry.AbsoluteExpiry.Compare(currentTime) == -1 || cacheEntry.VariableExpiry.Compare(currentTime) == -1) {
		deleteCacheEntry(cache, key)
		return nil, errors.New("key not found")
	}

	cacheEntry.VariableExpiry = cacheEntry.VariableExpiry.Add(cacheEntry.SlidingTimeout)
	writeMapEntry(cache, key, cacheEntry)
	return cacheEntry.Value, nil
}

func (cache *Cache) Add(key interface{}, value interface{}, config CacheEntryConfig) error {

	if _, ok := key.(*Cache); ok {
		return errors.New("the go_mem_cache.Cache type cannot be used as a key")
	}

	if _, ok := value.(*Cache); ok {
		return errors.New("the go_mem_cache.Cache type cannot be used as a value")
	}

	config.AbsoluteTimeout = config.AbsoluteTimeout.Abs()
	config.SlidingTimeout = config.SlidingTimeout.Abs()

	if cache.config.MaxEntries != -1 && len(cache.data) >= int(cache.config.MaxEntries) {
		return errors.New("maximum limit of cache entries reached")
	}

	if cache.config.MaxSizeBytes != -1 && reflect.TypeOf(cache.data).Size() >= uintptr(cache.config.MaxSizeBytes) {
		return errors.New("maximum limit of cache memory size reached")
	}

	var currentTime = time.Now().UTC()
	var cacheEntry cacheEntry
	cacheEntry.Value = value
	cacheEntry.AbsoluteExpiry = currentTime.Add(config.AbsoluteTimeout)
	cacheEntry.SlidingTimeout = config.SlidingTimeout
	cacheEntry.VariableExpiry = currentTime.Add(config.AbsoluteTimeout)
	cacheEntry.Expires = config.Expires
	writeMapEntry(cache, key, cacheEntry)
	return nil
}

func (cache *Cache) Remove(key interface{}) {
	deleteCacheEntry(cache, key)
}

func (cache *Cache) Clear() {
	for key := range cache.data {
		deleteCacheEntry(cache, key)
	}
}

func CleanCache(cache *Cache) {
	for {
		var currentTime = time.Now().UTC()

		for key := range cache.data {
			cacheEntry, _ := readMapEntry(cache, key)

			if cacheEntry.Expires && (cacheEntry.AbsoluteExpiry.Compare(currentTime) == -1 || cacheEntry.VariableExpiry.Compare(currentTime) == -1) {
				deleteCacheEntry(cache, key)
			}

		}

		time.Sleep(cache.config.ClearingCycleTime)
	}
}

func readMapEntry(cache *Cache, key interface{}) (cacheEntry, bool) {
	cache.RLock()
	defer cache.RUnlock()
	cacheEntry, ok := cache.data[key]
	return cacheEntry, ok
}

func writeMapEntry(cache *Cache, key interface{}, cacheEntry cacheEntry) {
	cache.Lock()
	defer cache.Unlock()
	cache.data[key] = cacheEntry
}

func deleteCacheEntry(cache *Cache, key interface{}) {
	cache.Lock()
	defer cache.Unlock()
	delete(cache.data, key)
}
