package ref

import (
	"errors"
	"mini_cache/asserts"
	"mini_cache/methods"
	"reflect"
)

type FunctionCache struct {
	MethodName string
	//ParamName  string
	FnValue  reflect.Value
	FnType   reflect.Type
	TypesIn  []reflect.Type
	TypesOut []reflect.Type
}

func ParseFn(fn interface{}) *FunctionCache {
	fc := FunctionCache{}
	fc.FnValue = reflect.ValueOf(fn)
	asserts.MustBeTrue(fc.FnValue.Kind() == reflect.Func, "must be function")
	fc.FnType = fc.FnValue.Type()
	fc.MethodName = methods.GetFuncName(fn)
	numTypesIn := fc.FnType.NumIn()
	fc.TypesIn = make([]reflect.Type, numTypesIn)
	for i := 0; i < numTypesIn; i++ {
		fc.TypesIn[i] = fc.FnType.In(i)
	}
	numTypesOut := fc.FnType.NumOut()
	fc.TypesOut = make([]reflect.Type, numTypesOut)
	for i := 0; i < numTypesOut; i++ {
		fc.TypesOut[i] = fc.FnType.Out(i)
	}
	return &fc
}

func (fc *FunctionCache) Call(params ...interface{}) []interface{} {
	pv := make([]reflect.Value, len(params))
	for i, p := range params {
		pv[i] = reflect.ValueOf(p)
	}
	rv := fc.FnValue.Call(pv)
	res := make([]interface{}, len(rv))
	for i, r := range rv {
		res[i] = r.Interface()
	}
	return res
}

type funcValidator func(cache *FunctionCache) error

// GenericType deprecated use Any instead
type GenericType int

type Any int

var (
	genericT = reflect.TypeOf(GenericType(1))
	anyT     = reflect.TypeOf(Any(1))
)

func NewFuncValidate(in []reflect.Type, out []reflect.Type) funcValidator {
	return func(cache *FunctionCache) error {
		isValid := func() bool {
			if in != nil {
				if len(in) != len(cache.TypesIn) {
					return false
				}
				for i, p := range in {
					if p != genericT && p != anyT && p != cache.TypesIn[i] {
						return false
					}
				}
			}
			if out != nil {
				if len(out) != len(cache.TypesOut) {
					return false
				}
				for i, p := range out {
					if p != genericT && p != anyT && p != cache.TypesOut[i] {
						return false
					}
				}
			}
			return true
		}
		if !isValid() {
			return errors.New("not match")
		}
		return nil
	}
}

func MakeTypeSlice(items ...interface{}) []reflect.Type {
	res := make([]reflect.Type, len(items))
	for i, each := range items {
		itemType := reflect.TypeOf(each)
		if itemType.Kind() == reflect.Ptr {
			itemType = itemType.Elem()
		}
		res[i] = itemType
	}
	return res
}
