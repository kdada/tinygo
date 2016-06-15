package web

import "github.com/kdada/tinygo/router"

// 高级执行器
type AdvancedExecutor struct {
	router.BaseRouterExecutor
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
		if this.ExecutePreFilters() {
			if this.StartMethod != nil {
				this.StartMethod.Call(context.Param)
			}
			var result = this.Method.Call(context.Param)
			if this.EndMethod != nil {
				this.EndMethod.Call(context.Param)
			}
			this.ExecutePostFilters(result)
			return result
		}
	}
	return nil
}

// ExecutePreFilters 执行全部前置过滤器
func (this *AdvancedExecutor) ExecutePreFilters() bool {
	var rec func(r router.Router, c router.RouterContext) bool
	rec = func(r router.Router, c router.RouterContext) bool {
		if r.Parent() != nil && !rec(r.Parent(), c) {
			return false
		}
		return r.ExecPreFilter(c)
	}
	return rec(this.End, this.Context)
}

// ExecutePostFilters 执行全部后置过滤器
func (this *AdvancedExecutor) ExecutePostFilters(result interface{}) bool {
	var r = this.End
	for r != nil {
		if !r.ExecPostFilter(this.RouterContext(), result) {
			return false
		}
		r = r.Parent()
	}
	return true
}
