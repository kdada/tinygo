package web

import (
	"reflect"
	"regexp"
	"strings"

	"github.com/kdada/tinygo/validator"
)

// 值提供器
type ValueProvider interface {
	// String 根据名称和类型返回相应的字符串值,返回的bool表示该值是否存在
	String(name string, t reflect.Type) ([]string, bool)
	// Value 根据名称和类型返回相应的对象
	Value(name string, t reflect.Type) interface{}
}

// 参数方法,根据name和type生成指定类型的数据
//type ParamFunc func(name string, t reflect.Type) interface{}

// 字段解析类型
type FieldKind byte

const (
	FieldKindUnimportant FieldKind = iota //次要字段 空 ,不验证直接注入值,无论是否成功都不报告错误
	FieldKindIgnore                       //忽略字段 -  ,忽略该字段,不验证也不注入,使用字段初始值
	FieldKindOptional                     //可选验证 ?  ,进行验证,验证通过后注入值,未通过验证使用字段初始值且不报告错误
	FieldKindRequired                     //必须验证 !  ,进行验证,验证通过后注入值,未通过验证报告错误
)

// 字段元数据
type FieldMetadata struct {
	Name      string              //字段名
	Field     reflect.StructField //字段信息
	Kind      FieldKind           //字段解析类型
	Validator validator.Validator //验证器
}

// Set 设置当前字段,instance必须是指针类型
func (this *FieldMetadata) Set(instance reflect.Value, param ValueProvider) error {
	if this.Kind != FieldKindIgnore {
		var valid = true
		if this.Kind != FieldKindUnimportant {
			// FieldKindOptional 和 FieldKindRequired 需要进行验证
			var strs []string
			strs, valid = param.String(this.Name, this.Field.Type)
			if this.Kind == FieldKindRequired && !valid {
				return ErrorRequiredField.Format(this.Name).Error()
			}
			//进行验证器验证,如果验证器不存在则默认为通过验证
			valid = true
			if this.Validator != nil {
				for i, v := range strs {
					valid = this.Validator.Validate(v)
					if !valid {
						if this.Kind == FieldKindRequired {
							if len(strs) == 1 {
								return ErrorFieldNotValid.Format(this.Name).Error()
							}
							return ErrorFieldsNotValid.Format(this.Name, i).Error()
						}
						break
					}
				}
			}
		}
		if valid {
			var f = param.Value(this.Name, this.Field.Type)
			var fValue reflect.Value
			if f == nil {
				fValue = reflect.New(this.Field.Type).Elem()
			} else {
				fValue = reflect.ValueOf(f)
			}
			instance.Elem().FieldByIndex(this.Field.Index).Set(fValue)
		}
	}
	return nil
}

// 结构体元数据
type StructMetadata struct {
	Name   string           //结构体名称
	Struct reflect.Type     //结构体类型
	Fields []*FieldMetadata //结构体字段元数据
}

// GenerateStruct 生成结构体,返回值为当前结构体的指针
func (this *StructMetadata) GenerateStruct(param ValueProvider) (reflect.Value, error) {
	var res = param.Value(this.Name, this.Struct)
	if res != nil {
		return reflect.ValueOf(res), nil
	}
	var result = reflect.New(this.Struct.Elem())
	for _, fMd := range this.Fields {
		var err = fMd.Set(result, param)
		if err != nil {
			return reflect.Value{}, err
		}
	}
	return result, nil
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
func (this *MethodMetadata) Call(param ValueProvider) ([]interface{}, error) {
	var params = make([]reflect.Value, 0)
	for _, sMd := range this.Params {
		var p, err = sMd.GenerateStruct(param)
		if err != nil {
			return nil, err
		}
		params = append(params, p)
	}
	var resultValue = this.Method.Call(params)
	var result = make([]interface{}, 0)
	for _, v := range resultValue {
		result = append(result, v.Interface())
	}
	return result, nil
}

// 全局结构体元数据信息
var globalStructMetadata = make(map[string]*StructMetadata)

// 验证字符串提取正则
var vldReg = regexp.MustCompile("^[?!] *?;(.*)$")

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
	ForeachField(s, func(field reflect.StructField) error {
		var tag = field.Tag.Get("vld")
		if tag == "" && field.Tag != "" {
			var firstChar = field.Tag[0]
			if firstChar == '!' || firstChar == '?' || firstChar == '-' {
				tag = string(field.Tag)
			}
		}
		var fMd = new(FieldMetadata)
		fMd.Name = field.Name
		fMd.Field = field
		switch {
		case tag == "":
			fMd.Kind = FieldKindUnimportant
		case strings.HasPrefix(tag, "!"):
			fMd.Kind = FieldKindRequired
		case strings.HasPrefix(tag, "?"):
			fMd.Kind = FieldKindOptional
		case strings.HasPrefix(tag, "-"):
			fMd.Kind = FieldKindIgnore
		default:
			return ErrorInvalidTag.Format(sMd.Name, field.Name, tag[0]).Error()
		}
		if fMd.Kind == FieldKindOptional || fMd.Kind == FieldKindRequired {
			//获取验证字符串
			var arr = vldReg.FindStringSubmatch(tag)
			if len(arr) == 2 {
				var src = arr[1]
				var vld, err = validator.NewValidator("string", src)
				if err != nil {
					return err
				}
				fMd.Validator = vld

			}
		}
		sMd.Fields = append(sMd.Fields, fMd)
		return nil
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
	ForeachResult(method.Type(), func(result reflect.Type) error {
		methodMd.Return = append(methodMd.Return, result)
		return nil
	})

	//遍历方法参数
	var err = ForeachParam(method.Type(), func(param reflect.Type) error {
		if !IsStructPtrType(param) {
			return ErrorParamNotPtr.Format(methodMd.Name, param.Kind().String()).Error()
		}
		var md, err = AnalyzeStruct(param)
		if err != nil {
			return err
		}
		methodMd.Params = append(methodMd.Params, md)
		return nil
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
func ForeachMethod(instance reflect.Type, solve func(method reflect.Method) error) error {
	if IsStructPtrType(instance) {
		for i := 0; i < instance.NumMethod(); i++ {
			var err = solve(instance.Method(i))
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// ForeachField 遍历instance的所有字段(包括匿名字段,但不包括私有字段),instance必须是结构体或结构体指针类型
func ForeachField(instance reflect.Type, solve func(field reflect.StructField) error) error {
	if IsStructPtrType(instance) {
		instance = instance.Elem()
	}
	if instance.Kind() == reflect.Struct {
		var preIndex = make([]int, 0, 2)
		return foreachField(preIndex, instance, solve)
	}
	return nil
}

// foreachField 遍历instance的所有字段(包括匿名字段,但不包括私有字段),instance必须是结构体或结构体指针类型
func foreachField(preIndex []int, instance reflect.Type, solve func(field reflect.StructField) error) error {
	for i := 0; i < instance.NumField(); i++ {
		var field = instance.Field(i)
		if (field.Name[0] >= 'A') && (field.Name[0] <= 'Z') {
			if field.Anonymous {
				preIndex = append(preIndex, i)
				foreachField(preIndex, field.Type, solve)
				preIndex = preIndex[:len(preIndex)-1]
			} else {
				var lenPI = len(preIndex)
				if lenPI > 0 {
					var index = make([]int, lenPI+len(field.Index))
					copy(index, preIndex)
					copy(index[lenPI:], field.Index)
					field.Index = index
				}
				var err = solve(field)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

// ForeachParam 遍历instance所有参数,instance必须是函数类型
func ForeachParam(instance reflect.Type, solve func(param reflect.Type) error) error {
	if instance.Kind() == reflect.Func {
		for i := 0; i < instance.NumIn(); i++ {
			var err = solve(instance.In(i))
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// ForeachResult 遍历instance所有返回值,instance必须是函数类型
func ForeachResult(instance reflect.Type, solve func(result reflect.Type) error) error {
	if instance.Kind() == reflect.Func {
		for i := 0; i < instance.NumOut(); i++ {
			var err = solve(instance.Out(i))
			if err != nil {
				return err
			}
		}
	}
	return nil
}
