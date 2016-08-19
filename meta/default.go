package meta

import "reflect"

// 全局值容器
var GlobalValueContainer = NewDefaultValueContainer()

// 生成方法
type GenerateFunc func() interface{}

// 默认值提供器,仅根据类型提供对象,不提供字符串值,因此无法用于字段校验
type DefaultValueProvider struct {
	Func GenerateFunc
}

// String 根据名称和类型返回相应的字符串值
func (this *DefaultValueProvider) String() []string {
	return []string{}
}

// Value 根据名称和类型返回相应的解析后的对象
func (this *DefaultValueProvider) Value() interface{} {
	return this.Func()
}

// 默认值容器
type DefaultValueContainer struct {
	NameContainer map[string]GenerateFunc //名称容器
	TypeContainer map[string]GenerateFunc //类型容器
}

// NewDefaultValueContainer 创建默认值容器
func NewDefaultValueContainer() *DefaultValueContainer {
	return &DefaultValueContainer{
		make(map[string]GenerateFunc),
		make(map[string]GenerateFunc),
	}
}

// Contains 检查值容器是否包含能够生成指定名称和类型的ValueProvider
func (this *DefaultValueContainer) Contains(name string, t reflect.Type) (ValueProvider, bool) {
	var f, ok = this.NameContainer[name]
	if ok {
		return &DefaultValueProvider{f}, ok
	}
	f, ok = this.TypeContainer[t.String()]
	if ok {
		return &DefaultValueProvider{f}, ok
	}
	return nil, false
}

// Register 注册指定名称或类型的生成器
//  name:
//    1.非空字符串:表示注册的是指定名称的生成器
//    2.非接口类型的值:表示注册的是指定非接口类型的生成器
//    3.值为nil的接口指针(例如:(*interface{})(nil)):表示注册的是指定接口类型的生成器
//    4.nil:自动按照generator参数判断应该生成哪种生成器
//  generator:
//    1.有一个返回值的函数:如果函数具有参数,则参数自动注入,在name为nil时将返回值作为生成器的类型
//    2.非接口类型的值(非指针):直接使用该值作为生成器的返回值,在name为nil时将值类型作为生成器的类型
//    3.非接口类型的指针(非nil):直接使用该值作为生成器的返回值,在name为nil时将指针类型作为生成器的类型
//    4.结构体的指针(为nil):则生成器自动创建该结构体的实例,并自动注入相应字段的值,在name为nil时将指针类型作为生成器的类型
func (this *DefaultValueContainer) Register(name interface{}, generator interface{}) error {
	if generator == nil {
		return ErrorMustBeStructPointer.Format("nil").Error()
	}
	var t = reflect.TypeOf(generator)
	if t.Kind() == reflect.Ptr && t.Elem().Kind() == reflect.Interface {
		//generator不能是接口指针
		return ErrorInvalidGenerator.Format(t.String()).Error()
	}
	var f GenerateFunc = nil
	var fName = ""
	var err error = nil
	if reflect.ValueOf(generator).IsNil() {
		// generator为nil时
		if !IsStructPtrType(t) {
			return ErrorMustBeStructPointer.Format(t.String()).Error()
		}
		//非接口类型的指针(为nil),不记录该类型的生成器,自动注入
		f, fName, err = this.pointer(generator)
	} else {
		if t.Kind() == reflect.Func {
			//函数
			f, fName, err = this.function(generator)
		} else {
			//非接口类型的值(非指针) 非接口类型的指针(为nil)
			f, fName, err = this.instance(generator)
		}
		if err != nil {
			return err
		}
	}
	var isStringName bool
	if name != nil {
		fName, isStringName = this.TranslateName(name)
	}
	//注册生成方法
	if isStringName {
		this.NameContainer[fName] = f
	} else {
		this.TypeContainer[fName] = f
	}
	return nil

}

// function 生成函数类型的GenerateFunc
//  generator:必须是有一个返回值的函数
//  return:(生成函数,返回值的类型全名,错误)
func (this *DefaultValueContainer) function(generator interface{}) (GenerateFunc, string, error) {
	var value = reflect.ValueOf(generator)
	var m, err = AnalyzeMethod("", &value)
	if err != nil {
		return nil, "", err
	}
	if len(m.Return) != 1 {
		return nil, "", ErrorInvalidFunction.Error()
	}
	return func() interface{} {
		var ins, err = m.Generate(this)
		if err != nil {
			return nil
		}
		var inss, ok = ins.([]interface{})
		if !ok || len(inss) != 1 {
			return nil
		}
		return inss[0]
	}, m.Return[0].String(), nil
}

// instance 生成函数类型的GenerateFunc
//  generator:必须是不为nil的值
//  return:(生成函数,类型全名,错误)
func (this *DefaultValueContainer) instance(generator interface{}) (GenerateFunc, string, error) {
	return func() interface{} {
		return generator
	}, reflect.TypeOf(generator).String(), nil
}

// pointer 生成结构体指针类型的GenerateFunc
//  generator:必须是结构体指针类型
//  return:(生成函数,返回值的类型全名,错误)
func (this *DefaultValueContainer) pointer(generator interface{}) (GenerateFunc, string, error) {
	var ptr = reflect.TypeOf(generator)
	var t = ptr.Elem()
	var g, err = AnalyzeStruct(t)
	if err != nil {
		return nil, "", err
	}
	var m = g.(*StructMetadata)
	return func() interface{} {
		var vp, ok = this.Contains(t.String(), t)
		if ok {
			//如果注册了generator的结构体类型,则使用相应的vp
			var v = vp.Value()
			var inst = reflect.New(t)
			inst.Elem().Set(reflect.ValueOf(v))
			return inst.Interface()
		} else {
			//没有注册相应结构体类型则自动创建
			var ins, err = m.New(this)
			if err != nil {
				return nil
			}
			return ins
		}
	}, ptr.String(), nil
}

// Delete 删除指定名称或类型的生成器
//  name:
//    1.非空字符串:表示注册的是指定名称的生成器
//    2.非接口类型的值:表示注册的是指定非接口类型的生成器
//    3.值为nil的接口指针(例如:(*interface{})(nil)):表示注册的是指定接口类型的生成器
func (this *DefaultValueContainer) Delete(name interface{}) {
	if name == nil {
		return
	}
	var trueName, isStringName = this.TranslateName(name)
	if isStringName {
		delete(this.NameContainer, trueName)
	} else {
		delete(this.TypeContainer, trueName)
	}
}

// TranslateName 将name转换为字符串类型
//  name:不能为nil
//  return:(名称或类型字符串,name是否是非空字符串)
func (this *DefaultValueContainer) TranslateName(name interface{}) (string, bool) {
	if name == nil {
		return "", false
	}
	var value, ok = name.(string)
	if value != "" && ok {
		return value, true
	}
	var t = reflect.TypeOf(name)
	if t.Kind() == reflect.Ptr && t.Elem().Kind() == reflect.Interface {
		return t.Elem().String(), false
	}
	return t.String(), false
}
