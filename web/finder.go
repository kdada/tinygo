package web

import (
	"reflect"
	"strconv"
	"time"
)

// 注册值查找器
func register(processor *HttpProcessor) {
	// 注册单类型查找器
	processor.RegisterFinder(reflect.TypeOf((*Context)(nil)), new(ContextCVF))
	processor.RegisterFinder(reflect.TypeOf((*FormFile)(nil)), new(FormFileCVF))
	processor.RegisterFinder(reflect.TypeOf(([]*FormFile)(nil)), new(FormFilesCVF))
	processor.RegisterFinder(reflect.TypeOf(time.Now()), new(TimeCVF))

	// 注册多类型查找器
	processor.RegisterMutiTypeFinder(new(MutiTypeCVF))
}

// http上下文值查找器
type ContextValueFinder interface {
	// Contains 查找context是否包含指定的值
	Contains(context *Context, name string, t reflect.Type) bool
	// String 查找context中指定的字符串值
	String(context *Context, name string, t reflect.Type) []string
	// Value 生成指定类型的值
	Value(context *Context, name string, t reflect.Type) interface{}
}

// *web.Context
type ContextCVF struct {
}

// Contains 查找context是否包含指定的值
func (this *ContextCVF) Contains(context *Context, name string, t reflect.Type) bool {
	return true
}

// String 查找context中指定的字符串值
func (this *ContextCVF) String(context *Context, name string, t reflect.Type) []string {
	return []string{}
}

// Value 生成指定类型的值
func (this *ContextCVF) Value(context *Context, name string, t reflect.Type) interface{} {
	return context
}

// *web.FormFile
type FormFileCVF struct {
}

// Contains 查找context是否包含指定的值
func (this *FormFileCVF) Contains(context *Context, name string, t reflect.Type) bool {
	var form = context.HttpContext.Request.MultipartForm
	if form != nil {
		var err = context.HttpContext.Request.ParseMultipartForm(int64(context.Processor.Config.MaxRequestMemory))
		if err != nil {
			return false
		}
		form = context.HttpContext.Request.MultipartForm
	}
	var _, ok = form.Value[name]
	return ok
}

// String 查找context中指定的字符串值
func (this *FormFileCVF) String(context *Context, name string, t reflect.Type) []string {
	return []string{}
}

// Value 生成指定类型的值
func (this *FormFileCVF) Value(context *Context, name string, t reflect.Type) interface{} {
	var f, err = context.ParamFile(name)
	if err == nil {
		return f
	}
	return nil
}

// []*web.FormFile
type FormFilesCVF struct {
	FormFileCVF
}

// Value 生成指定类型的值
func (this *FormFilesCVF) Value(context *Context, name string, t reflect.Type) interface{} {
	var f, err = context.ParamFiles(name)
	if err == nil {
		return f
	}
	return nil
}

// ValueCVF 单值查找器
type ValueCVF struct {
}

// Contains 查找context是否包含指定的值
func (this *ValueCVF) Contains(context *Context, name string, t reflect.Type) bool {
	var _, ok = context.Value(name)
	return ok
}

// String 查找context中指定的字符串值
func (this *ValueCVF) String(context *Context, name string, t reflect.Type) []string {
	var s, _ = context.Values(name)
	return s
}

// time.Time
type TimeCVF struct {
	ValueCVF
}

// Value 生成指定类型的值
func (this *TimeCVF) Value(context *Context, name string, t reflect.Type) interface{} {
	var v, ok = context.Value(name)
	if ok {
		//使用本地时区解析时间
		var result, err = time.ParseInLocation("2006-01-02 15:04:05", v, time.Local)
		if err == nil {
			return result
		}
	}
	return nil
}

// 多类型值查找器
type MutiTypeCVF struct {
}

// Contains 查找context是否包含指定的值
func (this *MutiTypeCVF) Contains(context *Context, name string, t reflect.Type) bool {
	var _, ok = context.Value(name)
	if !ok {
		return false
	}
	switch t.Kind() {
	case reflect.Bool, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Float32, reflect.Float64, reflect.String:
		return true
	case reflect.Slice:
		switch t.String() {
		case "[]string", "[]int", "[]float64", "[]bool":
			return true
		}
	}
	return false
}

// String 查找context中指定的字符串值
func (this *MutiTypeCVF) String(context *Context, name string, t reflect.Type) []string {
	var s, _ = context.Values(name)
	return s
}

// Value 生成指定类型的值
func (this *MutiTypeCVF) Value(context *Context, name string, t reflect.Type) interface{} {
	var value, ok = context.Values(name)
	if !ok {
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
