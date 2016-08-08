package web

import (
	"reflect"

	"github.com/kdada/tinygo/meta"
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
	Method *meta.MethodMetadata //执行方法
}

// NewAdvancedExecutor 创建高级执行器
func NewAdvancedExecutor(method *meta.MethodMetadata) *AdvancedExecutor {
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
			return this.Method.Generate(NewContextValueContainer(context))
		})
	}
	return nil, ErrorInvalidContext.Format(reflect.TypeOf(this.Context).String()).Error()

}
