package meta

import "reflect"

// 生成方法
type GenerateFunc func() interface{}

// 默认生成方法集合
var defaultGenerateFunc = make(map[string]GenerateFunc)

// RegisterNamedGenerateFunc 注册命名的生成方法
func RegisterNamedGenerateFunc(name string, gf GenerateFunc) {
	defaultGenerateFunc[name] = gf
}

// RegisterGenerateFunc 注册生成方法
func RegisterGenerateFunc(t reflect.Type, gf GenerateFunc) {
	RegisterNamedGenerateFunc(t.String(), gf)
}

// RegisterInterfaceType 注册接口类型,t类型必须能够赋值给i接口,并且t不能是接口类型
func RegisterInterfaceType(i reflect.Type, t reflect.Type) error {
	if t.Kind() != reflect.Interface {
		if t.AssignableTo(i) {
			if IsStructPtrType(t) {
				RegisterGenerateFunc(i, func() interface{} {
					return reflect.New(t.Elem()).Interface()
				})
			} else {
				RegisterGenerateFunc(i, func() interface{} {
					return reflect.New(t.Elem()).Elem().Interface()
				})
			}
			return nil
		}
		return ErrorNotAssignable.Format(t.String(), i.String()).Error()
	}
	return ErrorMustNotBeInterface.Error()
}

// RegisterType 注册类型,t必须是结构体或结构体指针
func RegisterType(t reflect.Type) error {
	return RegisterInterfaceType(t, t)
}

// RegisterInstance 注册实例类型
func RegisterInstance(ins interface{}) error {
	return RegisterType(reflect.TypeOf(ins))
}

// DefaultGenerateFuncByName 根据名称获取指定生成方法
func DefaultGenerateFuncByName(name string) (GenerateFunc, bool) {
	var f, ok = defaultGenerateFunc[name]
	return f, ok
}

// DefaultGenerateFunc 根据类型获取指定生成方法
func DefaultGenerateFunc(t reflect.Type) (GenerateFunc, bool) {
	return DefaultGenerateFuncByName(t.String())
}

// 默认值提供器,仅根据类型提供对象,不提供字符串值,因此无法用于字段校验
type DefaultValueProvider struct {
	f GenerateFunc
}

// String 根据名称和类型返回相应的字符串值
func (this *DefaultValueProvider) String() []string {
	return []string{}
}

// Value 根据名称和类型返回相应的解析后的对象
func (this *DefaultValueProvider) Value() interface{} {
	return this.f()
}

// 默认值容器
type DefaultValueContainer struct {
}

// NewDefaultValueContainer 创建默认值容器
func NewDefaultValueContainer() *DefaultValueContainer {
	return new(DefaultValueContainer)
}

// Contains 检查值容器是否包含能够生成指定名称和类型的ValueProvider
func (this *DefaultValueContainer) Contains(name string, t reflect.Type) (ValueProvider, bool) {
	var f, ok = DefaultGenerateFuncByName(name)
	if ok {
		return &DefaultValueProvider{f}, ok
	}
	f, ok = DefaultGenerateFunc(t)
	if ok {
		return &DefaultValueProvider{f}, ok
	}
	return nil, false
}
