package tests

import (
	"testing"

	gomemcache "github.com/Brokoth/go_mem_cache/src"
)

func newTestCache(t *testing.T) *gomemcache.Cache {
	t.Helper()
	cache, err := gomemcache.NewCache(gomemcache.DefaultCacheConfig())
	if err != nil {
		t.Fatalf("failed to create cache: %v", err)
	}
	return cache
}

func TestAddAndGetRoundTrip(t *testing.T) {
	cache := newTestCache(t)

	if err := cache.Add("key-1", "value-1", gomemcache.CacheEntryConfig{Expires: false}); err != nil {
		t.Fatalf("Add failed: %v", err)
	}

	value, err := cache.Get("key-1")
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}
	if value != "value-1" {
		t.Fatalf("expected %q, got %v", "value-1", value)
	}
}

func TestAddOverwritesExistingKey(t *testing.T) {
	cache := newTestCache(t)

	if err := cache.Add("key-1", "first", gomemcache.CacheEntryConfig{Expires: false}); err != nil {
		t.Fatalf("Add failed: %v", err)
	}
	if err := cache.Add("key-1", "second", gomemcache.CacheEntryConfig{Expires: false}); err != nil {
		t.Fatalf("Add failed: %v", err)
	}

	value, err := cache.Get("key-1")
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}
	if value != "second" {
		t.Fatalf("expected overwrite to update value, got %v", value)
	}
}

func TestGetMissingKeyReturnsError(t *testing.T) {
	cache := newTestCache(t)

	if _, err := cache.Get("does-not-exist"); err == nil {
		t.Fatalf("expected an error when getting a missing key")
	}
}

func TestCacheCannotBeUsedAsKey(t *testing.T) {
	cache := newTestCache(t)
	other := newTestCache(t)

	if err := cache.Add(other, "value", gomemcache.CacheEntryConfig{Expires: false}); err == nil {
		t.Fatalf("expected Add to reject a *Cache used as a key")
	}
}

func TestCacheCannotBeUsedAsValue(t *testing.T) {
	cache := newTestCache(t)
	other := newTestCache(t)

	if err := cache.Add("key", other, gomemcache.CacheEntryConfig{Expires: false}); err == nil {
		t.Fatalf("expected Add to reject a *Cache used as a value")
	}
}
