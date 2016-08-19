package meta

import "reflect"

// 方法元数据
type MethodMetadata struct {
	Name   string
	Method *reflect.Value
	Params []Generator
	Return []reflect.Type
}

// Generate 根据vc提供的值生成相应值
//  return:函数的返回值数组([]interface{})
func (this *MethodMetadata) Generate(vc ValueContainer) (interface{}, error) {
	var params = make([]reflect.Value, 0, len(this.Params))
	for _, sMd := range this.Params {
		var p, err = sMd.Generate(vc)
		if err != nil {
			return nil, err
		}
		params = append(params, reflect.ValueOf(p))
	}
	var resultValue = this.Method.Call(params)
	var result = make([]interface{}, 0, len(resultValue))
	for _, v := range resultValue {
		result = append(result, v.Interface())
	}
	return result, nil
}

// AnalyzeMethod 分析方法
//  name:方法名
//  method:方法值
func AnalyzeMethod(name string, method *reflect.Value) (*MethodMetadata, error) {
	if method.Type().Kind() != reflect.Func {
		return nil, ErrorParamNotFunc.Format(method.Type().String()).Error()
	}
	var mt = method.Type()
	var methodMd = new(MethodMetadata)
	methodMd.Name = name
	methodMd.Method = method
	methodMd.Params = make([]Generator, 0, mt.NumIn())
	methodMd.Return = make([]reflect.Type, 0, mt.NumOut())

	//遍历方法的返回值
	var err = ForeachResult(mt, func(result reflect.Type) error {
		methodMd.Return = append(methodMd.Return, result)
		return nil
	})
	if err != nil {
		return nil, err
	}
	//遍历方法参数
	err = ForeachParam(mt, func(param reflect.Type) error {
		var kind = param.Kind()
		var md Generator
		var e error
		if kind == reflect.Struct || IsStructPtrType(param) {
			md, e = AnalyzeStruct(param)
		} else {
			AnalyzeOther(param)
		}
		if e != nil {
			return e
		}
		methodMd.Params = append(methodMd.Params, md)
		return nil
	})
	if err != nil {
		return nil, err
	}
	return methodMd, nil
}

// AnalyzeStructMethod 分析结构体方法
func AnalyzeStructMethod(method *reflect.Method) (*MethodMetadata, error) {
	return AnalyzeMethod(method.Name, &method.Func)
}
