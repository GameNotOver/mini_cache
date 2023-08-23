package memorycache

import (
	"mini_cache/asserts"
	"mini_cache/cache"
	"time"
)

type Option func(*MemoryCache)

func OptCacheTTL(ttl time.Duration) Option {
	asserts.MustBeTrue(ttl >= time.Second, "ttl is too short")
	return func(c *MemoryCache) {
		c.ttl = ttl
	}
}

func OptName(name string) Option {
	asserts.MustBeTrue(name != "")
	return func(c *MemoryCache) {
		c.name = name
	}
}

// OptAttachTenantKey 是否租户隔离
func OptAttachTenantKey(isSeparate bool) Option {
	return func(c *MemoryCache) {
		c.attachTenantKey = isSeparate
	}
}

func OptWithLog(builder cache.LoggerBuilder) Option {
	return func(c *MemoryCache) {
		c.logBuilder = builder
	}
}
