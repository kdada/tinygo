package web

import (
	"reflect"
)

// IsStructPtrType 判断是否是结构体指针类型
func IsStructPtrType(instance reflect.Type) bool {
	return instance.Kind() == reflect.Ptr && instance.Elem().Kind() == reflect.Struct
}

// ForeachMethod 遍历instance的所有方法,instance必须是结构体指针类型(仅遍历使用结构体指针的方法)
func ForeachMethod(instance reflect.Type, solve func(method reflect.Method)) {
	if IsStructPtrType(instance) {
		for i := 0; i < instance.NumMethod(); i++ {
			solve(instance.Method(i))
		}
	}
}

// ForeachField 遍历instance的所有字段(包括匿名字段),instance必须是结构体或结构体指针类型
func ForeachField(instance reflect.Type, solve func(field reflect.StructField)) {
	if IsStructPtrType(instance) {
		instance = instance.Elem()
	}
	if instance.Kind() == reflect.Struct {
		for i := 0; i < instance.NumField(); i++ {
			var field = instance.Field(i)
			if field.Anonymous {
				ForeachField(field.Type, solve)
			} else {
				solve(field)
			}
		}
	}
}

// ForeachParam 遍历instance所有参数,instance必须是函数类型
func ForeachParam(instance reflect.Type, solve func(param reflect.Type)) {
	if instance.Kind() == reflect.Func {
		for i := 0; i < instance.NumIn(); i++ {
			solve(instance.In(i))
		}
	}
}

// ForeachResult 遍历instance所有返回值,instance必须是函数类型
func ForeachResult(instance reflect.Type, solve func(param reflect.Type)) {
	if instance.Kind() == reflect.Func {
		for i := 0; i < instance.NumOut(); i++ {
			solve(instance.Out(i))
		}
	}
}
