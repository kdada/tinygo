package web

import (
	"reflect"
	"strconv"
	"time"
)

// ParamContext 处理context类型
func ParamContext(context *Context, name string, t reflect.Type) interface{} {
	if t == reflect.TypeOf(context) {
		return context
	}
	return nil
}

// ParamTime 处理时间类型
func ParamTime(context *Context, name string, typ reflect.Type) interface{} {
	if name == "" {
		return time.Now()
	}
	var t, err = context.ParamString(name)
	if err == nil {
		//使用本地时区解析时间
		var result, err2 = time.ParseInLocation("2006-01-02 15:04:05", t, time.Local)
		if err2 == nil {
			return result
		}
	}
	return nil
}

// DefaultFunc 处理常规类型
func DefaultFunc(context *Context, name string, t reflect.Type) interface{} {
	if name == "" {
		return nil
	}
	var value, err = context.ParamStringArray(name)
	if err != nil {
		return nil
	}
	var first = ""
	if len(value) > 0 {
		first = value[0]
	}
	switch t.Kind() {
	case reflect.Bool:
		var result, err2 = strconv.ParseBool(first)
		if err2 == nil {
			return result
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		var result, err2 = strconv.ParseInt(first, 10, 64)
		if err2 == nil {
			switch t.Kind() {
			case reflect.Int:
				return int(result)
			case reflect.Int8:
				return int8(result)
			case reflect.Int16:
				return int16(result)
			case reflect.Int32:
				return int32(result)
			default:
				return result
			}
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		var result, err2 = strconv.ParseUint(first, 10, 64)
		if err2 == nil {
			switch t.Kind() {
			case reflect.Uint:
				return uint(result)
			case reflect.Uint8:
				return uint8(result)
			case reflect.Uint16:
				return uint16(result)
			case reflect.Uint32:
				return uint32(result)
			default:
				return result
			}
		}
	case reflect.Float32, reflect.Float64:
		var result, err2 = strconv.ParseFloat(first, 64)
		if err2 == nil {
			switch t.Kind() {
			case reflect.Float32:
				return float32(result)
			default:
				return result
			}
		}
	case reflect.Interface, reflect.String:
		return first
	case reflect.Slice:
		//可以解析的数组类型包括bool,int,float64,string类型
		switch t.String() {
		case "[]bool":
			boolValue := make([]bool, len(value), len(value))
			for i := 0; i < len(boolValue); i++ {
				boolValue[i], _ = strconv.ParseBool(value[i])
			}
			return boolValue
		case "[]int":
			intValue := make([]int, len(value), len(value))
			for i := 0; i < len(intValue); i++ {
				intValue[i], _ = strconv.Atoi(value[i])
			}
			return intValue
		case "[]float64":
			floatValue := make([]float64, len(value), len(value))
			for i := 0; i < len(floatValue); i++ {
				floatValue[i], _ = strconv.ParseFloat(value[i], 64)
			}
			return floatValue
		case "[]string":
			return value
		}
	}
	return nil
}

// 注册参数处理方法
func register(funcs map[string]ParamTypeFunc) {
	funcs[reflect.TypeOf((*Context)(nil)).String()] = ParamContext
	funcs[reflect.TypeOf(time.Now()).String()] = ParamContext
}
