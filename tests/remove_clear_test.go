package tests

import (
	"testing"

	gomemcache "github.com/Brokoth/go_mem_cache/src"
)

func TestRemoveDeletesEntry(t *testing.T) {
	cache := newTestCache(t)

	if err := cache.Add("key-1", "value-1", gomemcache.CacheEntryConfig{Expires: false}); err != nil {
		t.Fatalf("Add failed: %v", err)
	}

	cache.Remove("key-1")

	if _, err := cache.Get("key-1"); err == nil {
		t.Fatalf("expected Get to fail after Remove")
	}
}

func TestRemoveOnMissingKeyDoesNotPanic(t *testing.T) {
	cache := newTestCache(t)
	cache.Remove("does-not-exist")
}

func TestClearRemovesAllEntries(t *testing.T) {
	cache := newTestCache(t)

	keys := []string{"a", "b", "c"}
	for _, k := range keys {
		if err := cache.Add(k, k, gomemcache.CacheEntryConfig{Expires: false}); err != nil {
			t.Fatalf("Add(%q) failed: %v", k, err)
		}
	}

	cache.Clear()

	for _, k := range keys {
		if _, err := cache.Get(k); err == nil {
			t.Fatalf("expected Get(%q) to fail after Clear", k)
		}
	}
}
