package provider

import (
	"context"
	"github.com/ahmetb/go-linq/v3"
	"github.com/sirupsen/logrus"
	"mini_cache/asserts"
	"mini_cache/cache"
	"mini_cache/cache/memorycache"
	"mini_cache/cache/rediscache"
	"mini_cache/config"
)

type impl struct {
	conf   *config.Configs
	caches map[string]cache.Cache
}

func (i *impl) GetCache(id string) cache.Cache {
	c, ok := i.caches[id]
	asserts.MustBeTrue(ok, "cache (%s) is not configured", id)
	return c
}

func (i *impl) initMemoryCache(opt Option) cache.Cache {
	cacheOpts := []memorycache.Option{
		memorycache.OptCacheTTL(opt.GetTTL()),
		memorycache.OptName(opt.ID),
		memorycache.OptAttachTenantKey(opt.AttachTenantKey),
		memorycache.OptWithLog(func(ctx context.Context) cache.Logger {
			return logrus.WithContext(ctx)
		}),
	}
	size := 1000
	if opt.Size > 0 {
		size = opt.Size
	}
	c, err := memorycache.New(size, cacheOpts...)
	asserts.MustBeSuccess(err, "init memory cache fail: %s", err)
	return c
}

func (i *impl) initRedisCache(opt Option, conf *config.Configs) cache.Cache {
	opts := []rediscache.Option{
		rediscache.OptCacheTTL(opt.GetTTL()),
		rediscache.OptPrefix(opt.Prefix),
		rediscache.OptAttachTenantKey(opt.AttachTenantKey),
		rediscache.OptWithLog(func(ctx context.Context) cache.Logger {
			return logrus.WithContext(ctx)
		}),
	}
	if opt.BatchSize > 0 {
		opts = append(opts, rediscache.OptBatchSize(opt.BatchSize))
	}

	if opt.Compression {
		opts = append(opts, rediscache.OptSerializer(cache.CompressionSerializer))
	}
	redisOpt := conf.Redis["default"]
	redisProvider := rediscache.NewRedisProvider()
	return rediscache.New(redisProvider.NewRedis(&redisOpt), opts...)
}

func (i *impl) initCaches(options CacheOptions) {
	i.caches = make(map[string]cache.Cache)
	linq.From(options.Caches).ToMapBy(&i.caches, func(cacheOpt interface{}) interface{} {
		return cacheOpt.(Option).ID
	}, func(cacheOpt interface{}) interface{} {
		opt := cacheOpt.(Option)
		var c cache.Cache
		switch opt.Type {
		case "memory":
			c = i.initMemoryCache(opt)
		case "redis":
			c = i.initRedisCache(opt, i.conf)
		}
		return c
	})
}

func NewProviderFromConfig(
	conf *config.Configs,
	opts CacheOptions,
) CacheProvider {
	inst := &impl{
		conf: conf,
	}
	inst.initCaches(opts)
	return inst
}
