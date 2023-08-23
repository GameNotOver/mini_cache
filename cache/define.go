package cache

import (
	"context"
	"mini_cache/asserts"
	"mini_cache/ref"
	"reflect"
	"time"
)

type loadOptions struct {
	TTL        *time.Duration
	BatchSize  *int
	CacheSkipT *ref.FunctionCache
}

type Option func(options *loadOptions)

type Cache interface {
	BatchClear(ctx context.Context, keys []interface{}) error
	BatchLoadT(ctx context.Context, keys interface{}, loader interface{}, opts ...Option) (ResultSet, error)
}

func MergeOptions(options ...Option) *loadOptions {
	var dst loadOptions
	for _, each := range options {
		each(&dst)
	}
	return &dst
}

func BuildSkip(fc *ref.FunctionCache) func(interface{}) bool {
	if fc != nil {
		validate := ref.NewFuncValidate(
			ref.MakeTypeSlice(new(ref.GenericType)),
			ref.MakeTypeSlice(new(bool)),
		)
		asserts.MustBeSuccess(validate(fc))
	}
	return func(i interface{}) bool {
		return fc != nil && fc.Call(i)[0].(bool)
	}
}

func LoadTValidate(keys interface{}, loader interface{}) *ref.FunctionCache {
	keyVal := reflect.ValueOf(keys)
	asserts.MustBeTrue(keyVal.Kind() == reflect.Slice, "keys must be slice")
	asserts.MustBeTrue(keyVal.Len() > 0, "keys is empty")
	keyElemType := keyVal.Index(0).Type()
	fc := ref.ParseFn(loader)
	const msg = "loader must be like: func([]keyT) (map[keyT]valueT, error)"
	asserts.MustBeTrue(len(fc.TypesIn) == 1, msg)
	asserts.MustBeTrue(fc.TypesIn[0] == keyVal.Type(),
		"arguments type(%s) of loader must be %s", fc.TypesIn[0], keyVal.Type())
	asserts.MustBeTrue(len(fc.TypesOut) == 2, msg)
	errType := reflect.TypeOf((*error)(nil)).Elem()
	asserts.MustBeTrue(fc.TypesOut[1].Implements(errType), msg)
	asserts.MustBeTrue(fc.TypesOut[0].Kind() == reflect.Map, msg)
	asserts.MustBeTrue(fc.TypesOut[0].Key() == keyElemType, msg)
	return fc
}
