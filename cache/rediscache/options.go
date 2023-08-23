package rediscache

import (
	"mini_cache/asserts"
	"mini_cache/cache"
	"time"
)

// Options Redis 配置选项.
type Options struct {
	PSM            string `yaml:"psm" mapstructure:"psm"`
	Host           string `yaml:"host" mapstructure:"host"`
	Port           uint   `yaml:"port" mapstructure:"port"`
	Password       string `yaml:"password" mapstructure:"password"`
	ConnTimeoutMs  int64  `yaml:"conn_timeout" mapstructure:"conn_timeout"`
	ReadTimeoutMs  int64  `yaml:"read_timeout" mapstructure:"read_timeout"`
	WriteTimeoutMs int64  `yaml:"write_timeout" mapstructure:"write_timeout"`
	PoolTimeoutMs  int64  `yaml:"pool_timeout" mapstructure:"pool_timeout"`
	// Default is 10 connections per every CPU as reported by runtime.NumCPU
	PoolSize int `yaml:"pool_size"`
	// Default 10
	PoolInitSize int `yaml:"pool_init_size"`
	// 连接的 Redis 分库，默认 0
	Database int `yaml:"database"`
}

type Option func(*RedisCache)

func OptCacheTTL(ttl time.Duration) Option {
	asserts.MustBeTrue(ttl >= time.Second, "ttl is too short")
	return func(redisCache *RedisCache) {
		redisCache.ttl = ttl
	}
}

func OptPrefix(prefix string) Option {
	return func(redisCache *RedisCache) {
		redisCache.prefix = prefix
	}
}

func OptAttachTenantKey(isSeparate bool) Option {
	return func(redisCache *RedisCache) {
		redisCache.attachTenantKey = isSeparate
	}
}

func OptWithLog(builder cache.LoggerBuilder) Option {
	return func(redisCache *RedisCache) {
		redisCache.logBuilder = builder
	}
}

func OptBatchSize(size int) Option {
	asserts.MustBeTrue(size > 0, "size must be greater than 0")
	return func(redisCache *RedisCache) {
		redisCache.batchSize = size
	}
}

func OptSerializer(se cache.Serializer) Option {
	return func(redisCache *RedisCache) {
		redisCache.serializer = se
	}
}
