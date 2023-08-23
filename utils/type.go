package utils

import "reflect"

func IsSignedInteger(v reflect.Value) bool {
	kind := v.Kind()
	return kind >= reflect.Int && kind <= reflect.Int64
}

func IsUnsignedInteger(v reflect.Value) bool {
	kind := v.Kind()
	return kind >= reflect.Uint && kind <= reflect.Uint64
}

func IsInteger(v reflect.Value) bool {
	return IsSignedInteger(v) || IsUnsignedInteger(v)
}

func IsFloat(v reflect.Value) bool {
	kind := v.Kind()
	return kind == reflect.Float32 || kind == reflect.Float64
}

func IsNumber(v reflect.Value) bool {
	return IsInteger(v) || IsFloat(v)
}

func IsString(v reflect.Value) bool {
	return v.Kind() == reflect.String
}
