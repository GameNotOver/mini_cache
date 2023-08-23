package rediscache

import (
	"context"
	"fmt"
	"github.com/ahmetb/go-linq/v3"
	goredis "github.com/redis/go-redis/v9"
	"mini_cache/asserts"
	"mini_cache/cache"
	"mini_cache/cctx"
	"mini_cache/errors"
	"mini_cache/ref"
	"mini_cache/utils"
	"reflect"
	"time"
)

type modelGenerator func() reflect.Value

type RedisCache struct {
	prefix          string
	redisCli        *goredis.Client
	serializer      cache.Serializer
	ttl             time.Duration
	batchSize       int
	logBuilder      cache.LoggerBuilder
	skip            *ref.FunctionCache
	modelGen        modelGenerator
	attachTenantKey bool
}

func New(cli *goredis.Client, opt ...Option) *RedisCache {
	rc := &RedisCache{
		prefix:     "global",
		redisCli:   cli,
		serializer: cache.JsonSerializer,
		logBuilder: func(ctx context.Context) cache.Logger {
			return cache.EmptyLogger{}
		},
		ttl:       time.Minute * 10,
		batchSize: 200,
	}
	for _, each := range opt {
		each(rc)
	}
	return rc
}

type Cursor func() (idx int, key interface{}, value []byte, ok bool)

func (c *RedisCache) ensureKey(keys []interface{}) {
	for _, k := range keys {
		val := reflect.ValueOf(k)
		asserts.MustBeTrue(utils.IsNumber(val) || utils.IsString(val))
	}
}

func (c *RedisCache) getFormattedKeysWithTenantID(keys []interface{}, tenantID string) []string {
	formattedKeys := make([]string, 0, len(keys))
	for _, k := range keys {
		formattedKeys = append(formattedKeys, fmt.Sprintf("%s_%s:%v", c.prefix, tenantID, k))
	}
	return formattedKeys
}

func (c *RedisCache) getFormattedKeys(keys []interface{}) []string {
	formattedKeys := make([]string, 0, len(keys))
	for _, k := range keys {
		formattedKeys = append(formattedKeys, fmt.Sprintf("%s_%v", c.prefix, k))
	}
	return formattedKeys
}

func (c *RedisCache) formatKeys(ctx context.Context, keys []interface{}) []string {
	c.ensureKey(keys)
	tenantID := cctx.GetTenantIdFromContext(ctx)
	if c.attachTenantKey && tenantID != "" {
		return c.getFormattedKeysWithTenantID(keys, tenantID)
	} else {
		return c.getFormattedKeys(keys)
	}
}

func (c *RedisCache) redisCacheGenerator(ctx context.Context, keys []interface{}, batchSize int) (cursor Cursor) {
	logger := c.logBuilder(ctx)
	var curCursor int
	var curWindowBegin = -batchSize
	var curWindowEnd = 0
	keyLen := len(keys)
	formatKeys := c.formatKeys(ctx, keys)
	var batch [][]byte
	batchLoader := func(begin, end int) [][]byte {
		subKeys := formatKeys[begin:utils.Min(end, keyLen)]
		mResults, err := c.redisCli.MGet(ctx, subKeys...).Result()

		results := make([][]byte, len(subKeys))
		if err != nil && err != errors.Nil {
			logger.Warnf("Warning, fail to load cache from redis, %s, prefix: %s, offset/total: %d/%d", err.Error(),
				c.prefix, begin, len(keys))
			return results
		}

		for i, p := range mResults {
			if p == nil {
				continue
			}
			r, ok := p.(string)
			if ok && r != "" {
				results[i] = utils.StringToBytes(r)
			}
		}
		return results
	}
	loadFromCache := func() (idx int, key interface{}, value []byte, ok bool) {
		if curCursor >= curWindowEnd {
			return -1, nil, nil, false
		}
		ok = true
		idx = curCursor
		key = keys[curCursor]
		value = batch[curCursor-curWindowBegin]
		curCursor++
		return
	}
	return func() (idx int, key interface{}, value []byte, ok bool) {
		idx, key, value, ok = loadFromCache()
		if ok {
			return
		}
		// batch over, load from redis
		nextOffset := curWindowBegin + batchSize
		if nextOffset >= keyLen {
			return 0, nil, nil, false
		}
		curWindowBegin = nextOffset
		curWindowEnd = utils.Min(curWindowBegin+batchSize, keyLen)
		batch = batchLoader(curWindowBegin, curWindowEnd)
		return loadFromCache()
	}
}

func (c *RedisCache) batchSetCache(ctx context.Context, keys interface{}, data map[interface{}]interface{}, batchSize int, skip func(interface{}) bool, ttl time.Duration) {
	logger := c.logBuilder(ctx)
	type Pair struct {
		Key   string
		Value []byte
		Index int
	}
	linq.From(keys).Where(func(i interface{}) bool {
		v, found := data[i]
		return found && !skip(v)
	}).Select(func(i interface{}) interface{} {
		var p Pair
		p.Key = c.formatKeys(ctx, []interface{}{i})[0]
		if s, err := c.serializer.Serialize(data[i]); err != nil {
			logger.Warnf("[RedisCache][batchSetCache] WARNING, fail to serialize, %s", err.Error())
			return nil
		} else {
			p.Value = s
		}
		return &p
	}).SelectIndexed(func(i int, i2 interface{}) interface{} {
		p := i2.(*Pair)
		p.Index = i
		return p
	}).GroupBy(func(i interface{}) interface{} {
		return i.(*Pair).Index / batchSize
	}, func(i interface{}) interface{} {
		return i
	}).ForEachT(func(each linq.Group) {
		pipe := c.redisCli.Pipeline()
		linq.From(each.Group).ForEach(func(i interface{}) {
			p := i.(*Pair)
			if err := pipe.Set(ctx, p.Key, p.Value, ttl).Err(); err != nil {
				logger.Warnf("[RedisCache][batchSetCache] WARNING, fail to set, %s", err.Error())
			}
		})
		if pipe.Len() > 0 {
			if _, err := pipe.Exec(ctx); err != nil {
				logger.Warnf("[RedisCache][batchSetCache] WARNING, fail to exec pipeline, %s", err.Error())
			}
		}
	})
}

func (c *RedisCache) BatchClear(ctx context.Context, keys []interface{}) error {
	if len(keys) == 0 {
		return nil
	}
	formattedKeys := c.formatKeys(ctx, keys)
	rst := c.redisCli.Del(ctx, formattedKeys...)
	if rst.Err() != nil {
		return rst.Err()
	}
	return nil
}

func (c *RedisCache) BatchLoadT(ctx context.Context, keys interface{}, loader interface{}, opts ...cache.Option) (cache.ResultSet, error) {
	logger := c.logBuilder(ctx)
	fc := cache.LoadTValidate(keys, loader)
	valueType := fc.TypesOut[0].Elem()
	rs := cache.NewRsImpl(fc.TypesIn[0].Elem(), valueType)
	linq.From(keys).ToSlice(&rs.Keys)
	givenOptions := cache.MergeOptions(opts...)
	var batchSize = c.batchSize
	if givenOptions.BatchSize != nil {
		batchSize = *givenOptions.BatchSize
	}
	var (
		missing = reflect.MakeSlice(reflect.ValueOf(keys).Type(), 0, len(rs.Keys))
		hits    = make(map[interface{}]interface{}, len(rs.Keys))
	)
	cursor := c.redisCacheGenerator(ctx, rs.Keys, batchSize)
	for _, k, v, ok := cursor(); ok; _, k, v, ok = cursor() {
		if v == nil {
			missing = reflect.Append(missing, reflect.ValueOf(k))
			continue
		}
		val := reflect.New(valueType).Elem()
		if e := c.serializer.Deserialize(v, val.Addr().Interface()); e != nil {
			logger.Warnf("[RedisCache][BatchLoadT] fail to deserialize, %s", e.Error())
			missing = reflect.Append(missing, reflect.ValueOf(k))
			continue
		}
		hits[k] = val.Interface()
	}
	var loaderResults []interface{}
	var loaderGenerator *reflect.MapIter
	if missing.Len() != 0 {
		// from loader
		loaderResults = fc.Call(missing.Interface())
		if loaderResults[1] != nil {
			return nil, loaderResults[1].(error)
		}
		loaderGenerator = reflect.ValueOf(loaderResults[0]).MapRange()
		for loaderGenerator.Next() {
			k := loaderGenerator.Key().Interface()
			v := loaderGenerator.Value().Interface()
			hits[k] = v
		}
	}
	rs.Values = make([]interface{}, len(rs.Keys))
	for idx, k := range rs.Keys {
		if v, ok := hits[k]; ok {
			rs.Values[idx] = v
		} else {
			rs.Values[idx] = nil
		}
	}
	ttl := c.ttl
	if givenOptions.TTL != nil {
		ttl = *givenOptions.TTL
	}
	skip := cache.BuildSkip(c.skip)
	if givenOptions.CacheSkipT != nil {
		skip = cache.BuildSkip(givenOptions.CacheSkipT)
	}
	c.batchSetCache(ctx, missing.Interface(), hits, c.batchSize, skip, ttl)
	return rs, nil
}
