package meta

import "reflect"

// 值提供器
type ValueProvider interface {
	// Contains 检查值提供器是否包含指定类型的数据
	Contains(name string, t reflect.Type) bool
	// String 根据名称和类型返回相应的字符串值
	String(name string, t reflect.Type) []string
	// Value 根据名称和类型返回相应的解析后的对象
	Value(name string, t reflect.Type) interface{}
}

// 生成器
type Generator interface {
	// Generate 根据vp提供的值生成相应值
	Generate(vp ValueProvider) (interface{}, error)
}
