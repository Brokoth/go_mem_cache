package gomemcache

func New(config CacheConfig) Cache {
	var Cache Cache
	Cache.config = config
	return Cache
}

func Get(key string) {}

func Add(key string) {}

func Remove(key string) {}

func Clear() {}
