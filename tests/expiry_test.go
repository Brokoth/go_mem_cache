package tests

import (
	"testing"
	"time"

	gomemcache "github.com/Brokoth/go_mem_cache/src"
)

func TestNonExpiringEntryPersists(t *testing.T) {
	cache := newTestCache(t)

	cfg := gomemcache.CacheEntryConfig{
		Expires:         false,
		AbsoluteTimeout: 50 * time.Millisecond,
	}
	if err := cache.Add("key-1", "value-1", cfg); err != nil {
		t.Fatalf("Add failed: %v", err)
	}

	time.Sleep(150 * time.Millisecond)

	value, err := cache.Get("key-1")
	if err != nil {
		t.Fatalf("expected non-expiring entry to still be present, got error: %v", err)
	}
	if value != "value-1" {
		t.Fatalf("expected %q, got %v", "value-1", value)
	}
}

func TestEntryExpiresAfterAbsoluteTimeout(t *testing.T) {
	cache := newTestCache(t)

	cfg := gomemcache.CacheEntryConfig{
		Expires:         true,
		AbsoluteTimeout: 50 * time.Millisecond,
	}
	if err := cache.Add("key-1", "value-1", cfg); err != nil {
		t.Fatalf("Add failed: %v", err)
	}

	value, err := cache.Get("key-1")
	if err != nil {
		t.Fatalf("expected entry to be present before its timeout elapses, got error: %v", err)
	}
	if value != "value-1" {
		t.Fatalf("expected %q, got %v", "value-1", value)
	}

	time.Sleep(100 * time.Millisecond)

	if _, err := cache.Get("key-1"); err == nil {
		t.Fatalf("expected entry to have expired after its absolute timeout elapsed")
	}
}

// TestSlidingTimeoutDoesNotExtendPastAbsoluteTimeout documents the current
// behaviour of Get(): a successful Get extends the entry's VariableExpiry by
// SlidingTimeout, but expiry is triggered whenever *either* AbsoluteExpiry or
// VariableExpiry is in the past. Because AbsoluteExpiry is fixed at Add()
// time and is never itself extended, repeated Get calls do not keep an entry
// alive beyond its original AbsoluteTimeout.
func TestSlidingTimeoutDoesNotExtendPastAbsoluteTimeout(t *testing.T) {
	cache := newTestCache(t)

	cfg := gomemcache.CacheEntryConfig{
		Expires:         true,
		AbsoluteTimeout: 80 * time.Millisecond,
		SlidingTimeout:  5 * time.Second,
	}
	if err := cache.Add("key-1", "value-1", cfg); err != nil {
		t.Fatalf("Add failed: %v", err)
	}

	// Repeatedly access the entry well within the absolute timeout window.
	for i := 0; i < 3; i++ {
		if _, err := cache.Get("key-1"); err != nil {
			t.Fatalf("expected entry to be present on access %d, got error: %v", i, err)
		}
		time.Sleep(10 * time.Millisecond)
	}

	time.Sleep(100 * time.Millisecond)

	if _, err := cache.Get("key-1"); err == nil {
		t.Fatalf("expected entry to expire once its absolute timeout elapses, regardless of prior Get calls")
	}
}
