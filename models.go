package gomemcache

import (
	"time"
)

type CacheConfig struct {
	ClearingCycleTime time.Duration
	MaxEntries        int64
	MaxSizeBytes      int64
}

type Cache struct {
	data   map[interface{}]cacheEntry
	config CacheConfig
	Name   string
}

type cacheEntry struct {
	Value          interface{}
	SlidingTimeout time.Duration
	VariableExpiry time.Time
	AbsoluteExpiry time.Time
	Expires        bool
}

type CacheEntryConfig struct {
	SlidingTimeout  time.Duration
	AbsoluteTimeout time.Duration
	Expires         bool
}
