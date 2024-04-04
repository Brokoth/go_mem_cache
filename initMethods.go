package gomemcache

import (
	"errors"
	"time"
)

func NewCache(config CacheConfig) (resultCache *Cache, err error) {

	if config.MaxEntries < -1 || config.MaxEntries == 0 {
		return nil, errors.New("MaxEntries field in CacheConfig must be -1 or greater than 0")
	}

	if config.MaxSizeBytes < -1 || config.MaxSizeBytes == 0 {
		return nil, errors.New("MaxSizeBytes field in CacheConfig must be -1 or greater than 0")
	}

	config.ClearingCycleTime = config.ClearingCycleTime.Abs()
	var cache Cache
	cache.config = config
	go cache.CleanCache()
	return &cache, nil
}

func NewCacheEntryConfig() CacheEntryConfig {
	var cacheEntryConfig CacheEntryConfig
	cacheEntryConfig.SlidingTimeout = 0
	cacheEntryConfig.AbsoluteTimeout = 60 * time.Second
	cacheEntryConfig.Expires = true
	return cacheEntryConfig
}

func NewCacheConfig() CacheConfig {
	var cacheConfig CacheConfig
	cacheConfig.ClearingCycleTime = 60 * time.Second
	cacheConfig.MaxEntries = -1
	cacheConfig.MaxSizeBytes = -1
	return cacheConfig
}
