package web

import (
	"reflect"

	"github.com/kdada/tinygo/router"
)

// 公共执行器
type CommonExecutor struct {
	router.BaseRouterExecutor
}

// ExecutePreFilters 执行全部前置过滤器,返回结果决定了处理方法是否会被执行
func (this *CommonExecutor) ExecutePreFilters() bool {
	var rec func(r router.Router, c router.RouterContext) bool
	rec = func(r router.Router, c router.RouterContext) bool {
		if r.Parent() != nil && !rec(r.Parent(), c) {
			return false
		}
		return r.ExecPreFilter(c)
	}
	return rec(this.End, this.Context)
}

// ExecutePostFilters 执行全部后置过滤器,返回结果决定了处理方法的结果是否被Execute()返回
func (this *CommonExecutor) ExecutePostFilters(result interface{}) bool {
	var r = this.End
	for r != nil {
		if !r.ExecPostFilter(this.RouterContext(), result) {
			return false
		}
		r = r.Parent()
	}
	return true
}

// 简单执行器方法
type SimpleExecutorFunc func(r *Context) (interface{}, error)

// 简单执行器
type SimpleExecutor struct {
	CommonExecutor
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
		// 执行前置过滤器
		if this.ExecutePreFilters() {
			var r, err = this.f(context)
			if err != nil {
				return nil, err
			}
			//执行后置过滤器
			if this.ExecutePostFilters(r) {
				return r, nil
			}
		}
	}
	return nil, ErrorInvalidContext.Format(reflect.TypeOf(this.Context).String()).Error()
}

// 高级执行器
type AdvancedExecutor struct {
	CommonExecutor
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
		// 执行前置过滤器
		if this.ExecutePreFilters() {
			//执行处理方法
			var result, err = this.Method.Call(&ContextValueProvider{context})
			if err != nil {
				return nil, err
			}
			//执行后置过滤器
			if this.ExecutePostFilters(result) {
				return result, nil
			}
		}
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
