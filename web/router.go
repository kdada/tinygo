package web

import (
	"reflect"

	"github.com/kdada/tinygo/router"
)

// 创建适用于Web App的根路由
func NewRootRouter() (router.Router, error) {
	return router.NewRouter("base", "", nil)
}

// 创建控制器路由,instance为控制器对象
//  控制器方法必须满足如下格式:
//  func (this *SomeController) Method(param *ParamStruct) web.Result
//  this:必须是控制器指针
//  param:可以没有或者有多个,如果有则类型必须为结构体指针类型
//  返回结果必须为接口类型,并且必须能够赋值给web.Result接口
func NewControllerRouter(instance Controller) (router.Router, error) {
	var instanceType = reflect.TypeOf(instance)
	//var fields = make([]*FieldMetadata, 0)
	//遍历控制器字段
	ForeachField(instanceType, func(field reflect.StructField) {
		var meta = new(FieldMetadata)
		meta.Field = field
		meta.Name = field.Name
	})
	var resultType = reflect.TypeOf((*Result)(nil)).Elem()
	//var err error
	//遍历控制器方法
	ForeachMethod(instanceType, func(method reflect.Method) {
		//确认方法返回结果必须为1个
		if method.Type.NumOut() != 1 {
			return
		}
		//确认方法返回值为接口并且能够赋值给Result
		var result = method.Type.Out(0)
		if (result.Kind() != reflect.Interface) || !result.AssignableTo(resultType) {
			return
		}
		//遍历方法参数
		ForeachParam(method.Type, func(param reflect.Type) {
			if !IsStructPtrType(param) {
				//err =
			}
		})
	})
	return nil, nil
}
