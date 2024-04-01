package gomemcache

import (
	"time"
)

type CacheMethods interface {
	Get()
	Add()
	Remove()
	Clear()
}

type CacheConfig struct {
	ClearingCycleTime time.Duration
	MaxEntries        int64
	MaxSizeBytes      int64
}

type Cache struct {
	data   map[string]string
	config CacheConfig
	CacheMethods
}
