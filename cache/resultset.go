package cache

import (
	"github.com/ahmetb/go-linq/v3"
	"reflect"
)

type ResultSet interface {
	DecodeSliceT(dst interface{}) // dst: slice 的指针
	DecodeMapT(dst interface{})   // dst: map 的指针
}

type rsImpl struct {
	KeyType   reflect.Type
	ValueType reflect.Type
	Keys      []interface{}
	Values    []interface{}
}

func NewRsImpl(keyType, valueType reflect.Type) *rsImpl {
	return &rsImpl{
		KeyType:   keyType,
		ValueType: valueType,
	}
}

func (r *rsImpl) DecodeSliceT(dst interface{}) {
	results := reflect.MakeSlice(reflect.SliceOf(r.ValueType), 0, len(r.Keys))
	linq.From(r.Values).ForEach(func(i interface{}) {
		if i == nil {
			results = reflect.Append(results, reflect.Zero(r.ValueType))
		} else {
			results = reflect.Append(results, reflect.ValueOf(i))
		}
	})
	reflect.ValueOf(dst).Elem().Set(results)
}

func (r *rsImpl) DecodeMapT(dst interface{}) {
	type Item struct {
		Key   interface{}
		Value interface{}
	}
	source := make([]*Item, 0, len(r.Keys))
	linq.From(r.Keys).ForEachIndexed(func(idx int, key interface{}) {
		val := r.Values[idx]
		if val == nil {
			return
		}
		source = append(source, &Item{
			Key:   key,
			Value: val,
		})
	})
	// 先分配内存
	reflect.ValueOf(dst).Elem().
		Set(reflect.MakeMapWithSize(
			reflect.MapOf(r.KeyType, r.ValueType),
			len(r.Keys)))
	linq.From(source).ToMapBy(dst,
		func(each interface{}) interface{} { return each.(*Item).Key },
		func(each interface{}) interface{} { return each.(*Item).Value })
}
