package tests

import (
	"strconv"
	"sync"
	"testing"

	gomemcache "github.com/Brokoth/go_mem_cache/src"
)

func TestConcurrentAddAndGet(t *testing.T) {
	cache := newTestCache(t)

	const goroutines = 50
	var wg sync.WaitGroup
	wg.Add(goroutines)

	for i := 0; i < goroutines; i++ {
		go func(i int) {
			defer wg.Done()
			key := "key-" + strconv.Itoa(i)
			if err := cache.Add(key, i, gomemcache.CacheEntryConfig{Expires: false}); err != nil {
				t.Errorf("Add(%q) failed: %v", key, err)
				return
			}
			value, err := cache.Get(key)
			if err != nil {
				t.Errorf("Get(%q) failed: %v", key, err)
				return
			}
			if value != i {
				t.Errorf("expected %d, got %v", i, value)
			}
		}(i)
	}

	wg.Wait()
}
