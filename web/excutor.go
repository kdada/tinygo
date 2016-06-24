package web

import "github.com/kdada/tinygo/router"

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

// 简单执行器
type SimpleExecutor struct {
	CommonExecutor
	f func(r *Context) interface{}
}

// NewSimpleExecutor 创建简单执行器
func NewSimpleExecutor(f func(r *Context) interface{}) *SimpleExecutor {
	var se = new(SimpleExecutor)
	se.f = f
	return se
}

// Excute 执行
func (this *SimpleExecutor) Execute() interface{} {
	var context, ok = this.Context.(*Context)
	if ok {
		context.End = this.End
		// 执行前置过滤器
		if this.ExecutePreFilters() {
			var r = this.f(context)
			//执行后置过滤器
			if this.ExecutePostFilters(r) {
				return r
			}
		}
	}
	return nil
}

// 高级执行器
type AdvancedExecutor struct {
	CommonExecutor
	StartMethod *MethodMetadata //启动方法
	Method      *MethodMetadata //执行方法
	EndMethod   *MethodMetadata //结束方法
}

// NewAdvancedExecutor 创建高级执行器
func NewAdvancedExecutor(start, method, end *MethodMetadata) *AdvancedExecutor {
	var ae = new(AdvancedExecutor)
	ae.StartMethod = start
	ae.Method = method
	ae.EndMethod = end
	return ae
}

// Excute 执行
func (this *AdvancedExecutor) Execute() interface{} {
	var context, ok = this.Context.(*Context)
	if ok {
		context.End = this.End
		// 执行前置过滤器
		if this.ExecutePreFilters() {
			//执行Start方法
			if this.StartMethod != nil {
				this.StartMethod.Call(context.Param)
			}
			//执行处理方法
			var result = this.Method.Call(context.Param)
			//执行End方法
			if this.EndMethod != nil {
				this.EndMethod.Call(context.Param)
			}
			//执行后置过滤器
			if this.ExecutePostFilters(result) {
				return result
			}
		}
	}
	return nil
}
