package common

import (
	goredis "github.com/redis/go-redis/v9"
	"go.uber.org/dig"
	"mini_cache/cache/provider"
	"mini_cache/config"
)

type ProviderParams struct {
	dig.In
	Conf     *config.Configs
	RedisCli *goredis.Client
}

func NewProviderFromConfig(conf *config.Configs, opts provider.CacheOptions) provider.CacheProvider {
	redisProvider := provider.NewRedisProvider()
	inst := provider.NewProviderFromConfig(redisProvider, conf, opts)
	return inst
}
