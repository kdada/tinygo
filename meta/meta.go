// Package meta 包含元数据解析和验证工具
package meta

import "reflect"

// 值提供器
type ValueProvider interface {
	// String 根据名称和类型返回相应的字符串值,用于校验
	String() []string
	// Value 根据名称和类型返回相应的解析后的对象,用于注入
	Value() interface{}
}

// 值容器
type ValueContainer interface {
	// Contains 检查值容器是否包含能够生成指定名称和类型的ValueProvider
	Contains(name string, t reflect.Type) (ValueProvider, bool)
}

// 生成器
type Generator interface {
	// Generate 根据vc提供的值生成相应值
	Generate(vc ValueContainer) (interface{}, error)
}
