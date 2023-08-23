package provider

import (
	"mini_cache/cache"
)

type CacheProvider interface {
	GetCache(id string) cache.Cache
}
