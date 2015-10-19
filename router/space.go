package router

// 空间路由
// 空间路由仅用于隔离路由空间,本身并不具备任何功能
type SpaceRouter struct {
	BaseRouter
}

// NewSpaceRouter 创建空间路由
// name:路由名称
func NewSpaceRouter(name string) Router {
	var router = new(SpaceRouter)
	router.Init(name)
	return router
}

// Pass 传递指定的路由环境给当前的路由器
// context: 上下文环境
// return: 返回路由是否处理了该请求
// 如果请求已经被处理了,则该请求不应该继续被传递
func (this *SpaceRouter) Pass(context RouterContext) bool {
	var parts = context.RouterParts()
	if len(parts) > this.level {
		var pathName = parts[this.level]
		var childRouter, ok = this.Child(pathName)
		if ok {
			//执行前置过滤器
			ok = this.ExecBeforeFilter(context)
			if ok {
				//执行子路由处理方法
				ok = childRouter.Pass(context)
				//执行后置过滤器
				this.ExecAfterFilter(context)
				return ok
			}
		}
	}
	return false
}
