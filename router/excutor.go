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
	return this.FilterExecute(func() (interface{}, error) {
		return nil, ErrorExecutorDoNothing.Error()
	})
}

// FilterExecute 执行f方法并使用过滤器,未通过过滤器时返回相应错误
func (this *BaseRouterExecutor) FilterExecute(f func() (interface{}, error)) (interface{}, error) {
	// 执行前置过滤器
	if this.ExecutePreFilters() {
		//执行处理方法
		var result, err = f()
		if err != nil {
			return nil, err
		}
		//执行后置过滤器
		if this.ExecutePostFilters(result) {
			return result, nil
		}
		return nil, ErrorPostFilterNotPass.Error()
	}
	return nil, ErrorPreFilterNotPass.Error()
}

// ExecutePreFilters 执行全部前置过滤器,返回结果决定了处理方法是否会被执行
func (this *BaseRouterExecutor) ExecutePreFilters() bool {
	var rec func(r Router, c RouterContext) bool
	rec = func(r Router, c RouterContext) bool {
		if r.Parent() != nil && !rec(r.Parent(), c) {
			return false
		}
		return r.ExecPreFilter(c)
	}
	return rec(this.End, this.Context)
}

// ExecutePostFilters 执行全部后置过滤器,返回结果决定了处理方法的结果是否被Execute()返回
func (this *BaseRouterExecutor) ExecutePostFilters(result interface{}) bool {
	var r = this.End
	for r != nil {
		if !r.ExecPostFilter(this.Context, result) {
			return false
		}
		r = r.Parent()
	}
	return true
}
