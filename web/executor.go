package web

import (
	"reflect"

	"github.com/kdada/tinygo/router"
)

// 简单执行器方法
type SimpleExecutorFunc func(r *Context) (interface{}, error)

// 简单执行器
type SimpleExecutor struct {
	router.BaseRouterExecutor
	f SimpleExecutorFunc
}

// NewSimpleExecutor 创建简单执行器
func NewSimpleExecutor(f SimpleExecutorFunc) *SimpleExecutor {
	var se = new(SimpleExecutor)
	se.f = f
	return se
}

// Excute 执行
func (this *SimpleExecutor) Execute() (interface{}, error) {
	var context, ok = this.Context.(*Context)
	if ok {
		context.End = this.End
		return this.FilterExecute(func() (interface{}, error) {
			return this.f(context)
		})
	}
	return nil, ErrorInvalidContext.Format(reflect.TypeOf(this.Context).String()).Error()
}

// 高级执行器
type AdvancedExecutor struct {
	router.BaseRouterExecutor
	Method *MethodMetadata //执行方法
}

// NewAdvancedExecutor 创建高级执行器
func NewAdvancedExecutor(method *MethodMetadata) *AdvancedExecutor {
	var ae = new(AdvancedExecutor)
	ae.Method = method
	return ae
}

// Excute 执行
func (this *AdvancedExecutor) Execute() (interface{}, error) {
	var context, ok = this.Context.(*Context)
	if ok {
		context.End = this.End
		return this.FilterExecute(func() (interface{}, error) {
			return this.Method.Call(&ContextValueProvider{context})
		})
	}
	return nil, ErrorInvalidContext.Format(reflect.TypeOf(this.Context).String()).Error()

}

// http上下文值提供器
type ContextValueProvider struct {
	context *Context
}

// String 根据名称和类型返回相应的字符串值,返回的bool表示该值是否存在
func (this *ContextValueProvider) String(name string, t reflect.Type) ([]string, bool) {
	return this.context.Values(name)
}

// Value 根据名称和类型生成相应类型的数据,使用HttpProcessor中定义的参数生成方法
//  name:名称,根据该名称从Request里取值
//  t:生成类型,将对应值转换为该类型
//  return:返回指定类型的数据
func (this *ContextValueProvider) Value(name string, t reflect.Type) interface{} {
	var f = this.context.Processor.ParamFunc(t.String())
	if f != nil {
		return f(this.context, name, t)
	}
	return nil
}
