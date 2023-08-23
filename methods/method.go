package methods

import (
	"mini_cache/asserts"
	"reflect"
	"runtime"
	"strings"
)

func GetFuncName(fn interface{}) string {
	target := reflect.ValueOf(fn)
	asserts.MustBeTrue(target.Kind() == reflect.Func, "param must be function")
	fullName := runtime.FuncForPC(target.Pointer()).Name()
	splitList := strings.Split(fullName, ".")
	return splitList[len(splitList)-1]
}

func GetCallerName() string {
	pc, _, _, _ := runtime.Caller(2)
	fullName := runtime.FuncForPC(pc).Name()
	splitList := strings.Split(fullName, ".")
	return splitList[len(splitList)-1]
}

func GetCallerStructName() string {
	pc, _, _, _ := runtime.Caller(2)
	fullName := runtime.FuncForPC(pc).Name()
	splitList := strings.Split(fullName, ".")
	return splitList[len(splitList)-2]
}
