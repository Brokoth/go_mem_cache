package tests

import (
	"reflect"
	"testing"

	gomemcache "github.com/Brokoth/go_mem_cache/src"
)

func TestMaxEntriesLimitIsEnforced(t *testing.T) {
	cfg := gomemcache.DefaultCacheConfig()
	cfg.MaxEntries = 2
	cache, err := gomemcache.NewCache(cfg)
	if err != nil {
		t.Fatalf("failed to create cache: %v", err)
	}

	entryCfg := gomemcache.CacheEntryConfig{Expires: false}
	if err := cache.Add("key-1", "v1", entryCfg); err != nil {
		t.Fatalf("Add key-1 failed: %v", err)
	}
	if err := cache.Add("key-2", "v2", entryCfg); err != nil {
		t.Fatalf("Add key-2 failed: %v", err)
	}

	if err := cache.Add("key-3", "v3", entryCfg); err == nil {
		t.Fatalf("expected Add to fail once MaxEntries is reached")
	}

	// Freeing a slot should allow a subsequent Add to succeed again.
	cache.Remove("key-1")
	if err := cache.Add("key-3", "v3", entryCfg); err != nil {
		t.Fatalf("expected Add to succeed after freeing a slot, got error: %v", err)
	}
}

// mapHeaderSize returns the size, in bytes, that reflect.Type.Size() reports
// for any Go map type. A map value is always represented at runtime as a
// single pointer to an hmap struct, so this size is constant regardless of
// the map's key/value types or how many entries it holds.
func mapHeaderSize() int64 {
	return int64(reflect.TypeOf(map[string]string{}).Size())
}

// TestMaxSizeBytesLimitBehavior documents the current MaxSizeBytes
// enforcement in Add(): it compares MaxSizeBytes against
// reflect.TypeOf(cache.data).Size(), which is the fixed size of the map
// header type itself rather than the actual memory used by the cache's
// contents. As a result, once MaxSizeBytes clears that fixed size, the limit
// never blocks further Adds no matter how many entries are stored.
func TestMaxSizeBytesLimitBehavior(t *testing.T) {
	headerSize := mapHeaderSize()

	t.Run("rejects when MaxSizeBytes is at or below the map header size", func(t *testing.T) {
		cfg := gomemcache.DefaultCacheConfig()
		cfg.MaxSizeBytes = headerSize
		cache, err := gomemcache.NewCache(cfg)
		if err != nil {
			t.Fatalf("failed to create cache: %v", err)
		}

		if err := cache.Add("key-1", "v1", gomemcache.CacheEntryConfig{Expires: false}); err == nil {
			t.Fatalf("expected Add to fail when MaxSizeBytes <= the map header size")
		}
	})

	t.Run("never blocks once MaxSizeBytes exceeds the map header size", func(t *testing.T) {
		cfg := gomemcache.DefaultCacheConfig()
		cfg.MaxSizeBytes = headerSize + 1
		cache, err := gomemcache.NewCache(cfg)
		if err != nil {
			t.Fatalf("failed to create cache: %v", err)
		}

		for i := 0; i < 50; i++ {
			if err := cache.Add(i, i, gomemcache.CacheEntryConfig{Expires: false}); err != nil {
				t.Fatalf("expected Add #%d to succeed, got error: %v", i, err)
			}
		}
	})
}
