package memorycache

import (
	"context"
	"fmt"
	"github.com/ahmetb/go-linq/v3"
	"mini_cache/cache"
	"mini_cache/cctx"
	"reflect"
	"time"

	lru "github.com/hashicorp/golang-lru"
)

type MemoryCache struct {
	name            string
	ttl             time.Duration
	pool            *lru.TwoQueueCache
	logBuilder      cache.LoggerBuilder
	attachTenantKey bool
}

type entity struct {
	data     interface{}
	cachedAt time.Time
}

func New(size int, options ...Option) (*MemoryCache, error) {
	return newCache(size, options...)
}

func newCache(size int, options ...Option) (*MemoryCache, error) {
	queue, err := lru.New2Q(size)
	if err != nil {
		return nil, err
	}
	c := &MemoryCache{
		name: "default",
		ttl:  time.Minute * 10,
		logBuilder: func(ctx context.Context) cache.Logger {
			return cache.EmptyLogger{}
		},
		pool: queue,
	}
	for _, opt := range options {
		opt(c)
	}
	return c, nil
}

func formatKeyWithTenant(tenantID string, k interface{}) interface{} {
	return fmt.Sprintf("%s:%v", tenantID, k)
}

func (c *MemoryCache) formatRemove(ctx context.Context, k interface{}) {
	tenantID := cctx.GetTenantIdFromContext(ctx)
	if c.attachTenantKey && tenantID != "" {
		c.pool.Remove(formatKeyWithTenant(tenantID, k))
	} else {
		c.pool.Remove(k)
	}
}

func (c *MemoryCache) formatGet(ctx context.Context, k interface{}) (value interface{}, ok bool) {
	tenantID := cctx.GetTenantIdFromContext(ctx)
	if c.attachTenantKey && tenantID != "" {
		return c.pool.Get(formatKeyWithTenant(tenantID, k))
	}
	return c.pool.Get(k)
}

func (c *MemoryCache) formatAdd(ctx context.Context, k, v interface{}) {
	tenantID := cctx.GetTenantIdFromContext(ctx)
	if c.attachTenantKey && tenantID != "" {
		c.pool.Add(formatKeyWithTenant(tenantID, k), v)
	} else {
		c.pool.Add(k, v)
	}
}

func (c *MemoryCache) BatchClear(ctx context.Context, keys []interface{}) error {
	for _, k := range keys {
		c.formatRemove(ctx, k)
	}
	return nil
}

func (c *MemoryCache) BatchLoadT(ctx context.Context, keys interface{}, loader interface{}, opts ...cache.Option) (cache.ResultSet, error) {
	fc := cache.LoadTValidate(keys, loader)
	valueType := fc.TypesOut[0].Elem()
	rs := cache.NewRsImpl(fc.TypesIn[0].Elem(), valueType)
	linq.From(keys).ToSlice(&rs.Keys)
	givenOptions := cache.MergeOptions(opts...)
	var (
		missing = reflect.MakeSlice(reflect.ValueOf(keys).Type(), 0, len(rs.Keys))
		hits    = make(map[interface{}]interface{}, len(rs.Keys))
	)
	ttl := c.ttl
	if givenOptions.TTL != nil {
		ttl = *givenOptions.TTL
	}
	for _, k := range rs.Keys {
		ent, exists := c.formatGet(ctx, k)
		if !exists {
			missing = reflect.Append(missing, reflect.ValueOf(k))
			continue
		}
		wrapped := ent.(*entity)
		if ttl != 0 && time.Now().Sub(wrapped.cachedAt) > ttl {
			missing = reflect.Append(missing, reflect.ValueOf(k))
			continue
		}
		hits[k] = wrapped.data
	}
	var loaderResults []interface{}
	var loaderGenerator *reflect.MapIter
	if missing.Len() != 0 {
		// from loader 函数调用
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
	now := time.Now()
	linq.From(missing.Interface()).ForEach(func(i interface{}) {
		v, found := hits[i]
		if !found {
			return
		}
		c.formatAdd(ctx, i, &entity{
			data:     v,
			cachedAt: now,
		})
	})
	return rs, nil
}
