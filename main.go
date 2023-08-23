package main

import (
	"context"
	"fmt"
	"mini_cache/cache"
	"mini_cache/cache/provider"
	"mini_cache/toolbox"
)

func CacheFunc(missKeys []string) (map[string]string, error) {
	dataMap := map[string]string{
		"key_1": "val_1",
	}
	resMap := make(map[string]string, len(missKeys))
	for _, key := range missKeys {
		resMap[key] = dataMap[key]
	}
	return resMap, nil
}

func UseRedisCache(ctx context.Context, instance provider.CacheProvider) map[string]string {
	resMap := make(map[string]string, 2)
	redisCache := instance.GetCache("cache:redis")
	rs, _ := redisCache.BatchLoadT(ctx, []string{"1", "2"}, func(missKeys []string) (r map[string]string, _ error) {
		return map[string]string{
			"1": "3",
			"2": "4",
		}, nil
	})
	rs.DecodeMapT(&resMap)
	return resMap
}

func UseMemoryCache(ctx context.Context, instance provider.CacheProvider) map[string]string {
	resMap := make(map[string]string, 2)
	redisCache := instance.GetCache("cache:memory")
	rs, _ := redisCache.BatchLoadT(ctx, []string{"1", "2"}, func(missKeys []string) (r map[string]string, _ error) {
		return map[string]string{
			"1": "3",
			"2": "4",
		}, nil
	})
	rs.DecodeMapT(&resMap)
	return resMap
}

func main() {
	var ctx = context.Background()
	//initialize.Init()
	conf := toolbox.LoadConfig()
	opts := toolbox.LoadOpts()
	instance := provider.NewProviderFromConfig(conf, opts)
	resMap := UseRedisCache(ctx, instance)
	resMap2 := UseMemoryCache(ctx, instance)
	resMap2 = UseMemoryCache(ctx, instance)
	keys := []string{"1", "2"}
	cache.LoadTValidate(keys, CacheFunc)
	fmt.Printf("hello: %v, %v", resMap, resMap2)
}
