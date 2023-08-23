package rediscache

import (
	"fmt"
	goredis "github.com/redis/go-redis/v9"
	"strconv"
)

type RedisProvider interface {
	NewRedis(conf *Options) *goredis.Client
}

type RedisProviderImpl struct{}

func NewRedisProvider() RedisProvider {
	return &RedisProviderImpl{}
}

func (r *RedisProviderImpl) NewRedis(conf *Options) *goredis.Client {
	opts := &goredis.Options{}
	opts.DB = conf.Database
	if conf.PoolInitSize > 0 {
		opts.PoolSize = conf.PoolInitSize
	}
	if conf.Host != "" && conf.Port != 0 {
		opts.Addr = conf.Host + ":" + strconv.Itoa(int(conf.Port))
	}
	if conf.Password != "" {
		opts.Password = conf.Password
	}
	fmt.Printf("%s:%d", conf.Host, conf.Port)
	return goredis.NewClient(opts)
}
