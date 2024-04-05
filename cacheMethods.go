package gomemcache

import (
	"errors"
	"log"
	"reflect"
	"time"
)

func (cache Cache) Get(key interface{}) (value interface{}, err error) {
	var currentTime = time.Now().UTC()
	cacheEntry, ok := cache.data[key]

	if !ok {
		return nil, errors.New("key not found")
	}

	if cacheEntry.Expires && (cacheEntry.AbsoluteExpiry.Compare(currentTime) == -1 || cacheEntry.VariableExpiry.Compare(currentTime) == -1) {
		delete(cache.data, key)
		return nil, errors.New("key not found")
	}

	cacheEntry.VariableExpiry = cacheEntry.VariableExpiry.Add(cacheEntry.SlidingTimeout)
	return cacheEntry.Value, nil
}

func (cache Cache) Add(key interface{}, value interface{}, config CacheEntryConfig) error {
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
	cache.data[key] = cacheEntry
	return nil
}

func (cache Cache) Remove(key interface{}) {
	delete(cache.data, key)
}

func (cache Cache) Clear() {
	for k := range cache.data {
		delete(cache.data, k)
	}
}

func (cache Cache) CleanCache() {
	for {
		log.Println("Size of cache before cleaning: ", len(cache.data))
		var currentTime = time.Now().UTC()

		for key := range cache.data {
			cacheEntry := cache.data[key]

			if cacheEntry.Expires && (cacheEntry.AbsoluteExpiry.Compare(currentTime) == -1 || cacheEntry.VariableExpiry.Compare(currentTime) == -1) {
				delete(cache.data, key)
			}

		}
		log.Println("Size of cache after cleaning: ", len(cache.data))
		time.Sleep(cache.config.ClearingCycleTime)
	}
}
