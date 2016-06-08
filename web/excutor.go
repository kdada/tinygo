package web

import (
	"reflect"

	"github.com/kdada/tinygo/router"
)

// MethodExecutor
type MethodExecutor struct {
	router.BaseRouterExecutor
	Method func(end router.Router, context router.RouterContext) interface{}
}

// Excute 执行
func (this *MethodExecutor) Execute() interface{} {
	return this.Method(this.End, this.Context)
}

// 字段元数据
type FieldMetadata struct {
	Name  string              //字段名
	Field reflect.StructField //字段信息
}

// 参数元数据
type ParamMetadata struct {
	Type   reflect.Type     //参数类型
	Fields []*FieldMetadata //字段元数据
}

// 高级执行器
type AdvancedExecutor struct {
	router.BaseRouterExecutor
	Instance reflect.Type     //控制器类型
	Fields   []*FieldMetadata //控制器字段元数据(router.Router,router.RouterContext类型的公开字段将被自动设置)
	Method   *reflect.Value   //方法值,方法参数必须为结构体,方法返回值必须符合Result接口
	In       []*ParamMetadata //参数元数据
}

// Excute 执行
func (this *AdvancedExecutor) Execute() interface{} {
	return nil
}
