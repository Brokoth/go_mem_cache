package tests

import (
	"testing"
	"time"

	gomemcache "github.com/Brokoth/go_mem_cache/src"
)

func TestDefaultCacheConfig(t *testing.T) {
	cfg := gomemcache.DefaultCacheConfig()

	if cfg.ClearingCycleTime != time.Minute {
		t.Fatalf("expected ClearingCycleTime of 1 minute, got %v", cfg.ClearingCycleTime)
	}
	if cfg.MaxEntries != -1 {
		t.Fatalf("expected unlimited MaxEntries (-1), got %d", cfg.MaxEntries)
	}
	if cfg.MaxSizeBytes != -1 {
		t.Fatalf("expected unlimited MaxSizeBytes (-1), got %d", cfg.MaxSizeBytes)
	}
}

func TestDefaultCacheEntryConfig(t *testing.T) {
	cfg := gomemcache.DefaultCacheEntryConfig()

	if cfg.SlidingTimeout != 0 {
		t.Fatalf("expected SlidingTimeout of 0, got %v", cfg.SlidingTimeout)
	}
	if cfg.AbsoluteTimeout != 60*time.Second {
		t.Fatalf("expected AbsoluteTimeout of 60s, got %v", cfg.AbsoluteTimeout)
	}
	if !cfg.Expires {
		t.Fatalf("expected Expires to be true")
	}
}

func TestNewCacheAcceptsValidConfig(t *testing.T) {
	configs := []gomemcache.CacheConfig{
		gomemcache.DefaultCacheConfig(),
		{ClearingCycleTime: time.Second, MaxEntries: -1, MaxSizeBytes: -1},
		{ClearingCycleTime: time.Second, MaxEntries: 5, MaxSizeBytes: 1024},
	}

	for _, cfg := range configs {
		cache, err := gomemcache.NewCache(cfg)
		if err != nil {
			t.Fatalf("expected NewCache(%+v) to succeed, got error: %v", cfg, err)
		}
		if cache == nil {
			t.Fatalf("expected a non-nil cache for config %+v", cfg)
		}
	}
}

func TestNewCacheRejectsInvalidMaxEntries(t *testing.T) {
	invalidValues := []int64{0, -2, -100}

	for _, v := range invalidValues {
		cfg := gomemcache.DefaultCacheConfig()
		cfg.MaxEntries = v
		if _, err := gomemcache.NewCache(cfg); err == nil {
			t.Fatalf("expected NewCache to reject MaxEntries=%d", v)
		}
	}
}

func TestNewCacheRejectsInvalidMaxSizeBytes(t *testing.T) {
	invalidValues := []int64{0, -2, -100}

	for _, v := range invalidValues {
		cfg := gomemcache.DefaultCacheConfig()
		cfg.MaxSizeBytes = v
		if _, err := gomemcache.NewCache(cfg); err == nil {
			t.Fatalf("expected NewCache to reject MaxSizeBytes=%d", v)
		}
	}
}
