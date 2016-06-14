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
		this.Method.Call(context.Param)
	}
	return nil
}
