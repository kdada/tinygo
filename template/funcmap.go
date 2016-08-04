package template

import (
	"html/template"
	"reflect"
)

// 公共模版方法
var commonFuncMap template.FuncMap

// RegisterTemplateFunc 注册模板函数
//  f:必须是一个函数,并且只能有一个返回值,或者有两个返回值并且第二个返回值为error
func RegisterTemplateFunc(name string, f interface{}) error {
	if f == nil {
		return ErrorParamMustBeFunc.Error()
	}
	var t = reflect.TypeOf(f)
	if t.Kind() == reflect.Func {
		commonFuncMap[name] = f
		return nil
	}
	return ErrorParamMustBeFunc.Error()
}

// DeleteTemplateFunc 删除模板函数
func DeleteTemplateFunc(name string) error {
	delete(commonFuncMap, name)
	return nil
}
