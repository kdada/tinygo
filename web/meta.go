package web

import "reflect"

// 参数方法,根据name和type生成指定类型的数据
type ParamFunc func(name string, t reflect.Type) interface{}

// 字段元数据
type FieldMetadata struct {
	Name  string              //字段名
	Field reflect.StructField //字段信息
}

// Set 设置当前字段,instance必须是指针类型
func (this *FieldMetadata) Set(instance reflect.Value, param ParamFunc) {
	var f = param(this.Name, this.Field.Type)
	var fValue reflect.Value
	if f == nil {
		fValue = reflect.New(this.Field.Type).Elem()
	} else {
		fValue = reflect.ValueOf(f)
	}
	instance.Elem().FieldByIndex(this.Field.Index).Set(fValue)
}

// 结构体元数据
type StructMetadata struct {
	Name   string           //结构体名称
	Struct reflect.Type     //结构体类型
	Fields []*FieldMetadata //结构体字段元数据
}

// GenerateStruct 生成结构体,返回值为当前结构体的指针
func (this *StructMetadata) GenerateStruct(param ParamFunc) reflect.Value {
	var res = param(this.Name, this.Struct)
	if res != nil {
		return reflect.ValueOf(res)
	}
	var result = reflect.New(this.Struct.Elem())
	for _, fMd := range this.Fields {
		fMd.Set(result, param)
	}
	return result
}

// 方法元数据
type MethodMetadata struct {
	Name   string
	Method reflect.Value
	Params []*StructMetadata
	Return []reflect.Type
}

// Call 调用当前方法元数据中包含的方法
//  param:通过字段名和类型名返回相应的值,返回nil表示该字段需要自动设置为默认值
//  return:返回方法的返回值
func (this *MethodMetadata) Call(param ParamFunc) []interface{} {
	var params = make([]reflect.Value, 0)
	for _, sMd := range this.Params {
		params = append(params, sMd.GenerateStruct(param))
	}
	var resultValue = this.Method.Call(params)
	var result = make([]interface{}, 0)
	for _, v := range resultValue {
		result = append(result, v.Interface())
	}
	return result
}

// 全局结构体元数据信息
var globalStructMetadata = make(map[string]*StructMetadata)

// AnalyzeStruct 分析结构体字段(包括匿名字段)
func AnalyzeStruct(s reflect.Type) (*StructMetadata, error) {
	if !IsStructPtrType(s) {
		return nil, ErrorNotStructPtr.Format(s.String()).Error()
	}
	var sm, ok = globalStructMetadata[s.String()]
	if ok {
		return sm, nil
	}
	var sMd = new(StructMetadata)
	sMd.Name = s.Name()
	sMd.Struct = s
	sMd.Fields = make([]*FieldMetadata, 0)
	ForeachField(s, func(field reflect.StructField) {
		var fMd = new(FieldMetadata)
		fMd.Name = field.Name
		fMd.Field = field
		sMd.Fields = append(sMd.Fields, fMd)
	})
	globalStructMetadata[s.String()] = sMd
	return sMd, nil
}

// AnalyzeMethod 分析方法(方法参数必须为结构体指针类型)
func AnalyzeMethod(name string, method reflect.Value) (*MethodMetadata, error) {
	if method.Type().Kind() != reflect.Func {
		return nil, ErrorNotMethod.Format(method.Type()).Error()
	}
	var methodMd = new(MethodMetadata)
	methodMd.Name = name
	methodMd.Method = method
	methodMd.Params = make([]*StructMetadata, 0)
	methodMd.Return = make([]reflect.Type, 0)

	//遍历方法的返回值
	ForeachResult(method.Type(), func(result reflect.Type) {
		methodMd.Return = append(methodMd.Return, result)
	})

	var err error
	//遍历方法参数
	ForeachParam(method.Type(), func(param reflect.Type) {
		if !IsStructPtrType(param) {
			err = ErrorParamNotPtr.Format(methodMd.Name, param.Kind().String()).Error()
			return
		}
		var md, err2 = AnalyzeStruct(param)
		if err2 != nil {
			err = err2
			return
		}
		methodMd.Params = append(methodMd.Params, md)
	})
	if err != nil {
		return nil, err
	}
	return methodMd, nil
}

// AnalyzeControllerMethod 分析控制器方法(方法参数必须为结构体指针类型)
func AnalyzeControllerMethod(method reflect.Method) (*MethodMetadata, error) {
	return AnalyzeMethod(method.Name, method.Func)
}

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

// ForeachField 遍历instance的所有字段(包括匿名字段,但不包括私有字段),instance必须是结构体或结构体指针类型
func ForeachField(instance reflect.Type, solve func(field reflect.StructField)) {
	if IsStructPtrType(instance) {
		instance = instance.Elem()
	}
	if instance.Kind() == reflect.Struct {
		for i := 0; i < instance.NumField(); i++ {
			var field = instance.Field(i)
			if (field.Name[0] >= 'A') && (field.Name[0] <= 'Z') {
				if field.Anonymous {
					ForeachField(field.Type, solve)
				} else {
					solve(field)
				}
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
func ForeachResult(instance reflect.Type, solve func(result reflect.Type)) {
	if instance.Kind() == reflect.Func {
		for i := 0; i < instance.NumOut(); i++ {
			solve(instance.Out(i))
		}
	}
}
