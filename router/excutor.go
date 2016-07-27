package router

//基础路由执行器
type BaseRouterExecutor struct {
	End     Router
	Context RouterContext
}

// Router 返回生成RouterExcutor的路由
func (this *BaseRouterExecutor) Router() Router {
	return this.End
}

// SetRouter 设置生成RouterExcutor的路由
func (this *BaseRouterExecutor) SetRouter(router Router) {
	this.End = router
}

// RouterContext 返回路由上下文
func (this *BaseRouterExecutor) RouterContext() RouterContext {
	return this.Context
}

// SetRouterContext 设置路由上下文
func (this *BaseRouterExecutor) SetRouterContext(context RouterContext) {
	this.Context = context
}

// Excute 执行
func (this *BaseRouterExecutor) Execute() (interface{}, error) {
	return nil, ErrorExecutorDoNothing.Error()
}
