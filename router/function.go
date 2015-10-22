package router

// 函数路由
type FunctionRouter struct {
	BaseRouter
	function func(RouterContext) //执行函数
}

// Pass 传递指定的路由环境给当前的路由器
// context: 上下文环境
// return: 返回路由是否处理了该请求
// 如果请求已经被处理了,则该请求不应该继续被传递
func (this *FunctionRouter) Pass(context RouterContext) bool {
	var parts = context.RouterParts()
	if len(parts) == this.Level()+1 {
		//检查当前路由是否能够处理
		var route = parts[this.Level()]
		var routeData, ok = this.check(route)
		if ok {
			//添加路由参数
			for k, v := range routeData {
				context.AddRouterParams(k, v)
			}
			//处理
			var fre = &FunctionRouterExecutor{this}
			context.AddRouter(this)
			context.AddContextExecutor(fre)
			return true
		}
	}

	return false
}

// Restful控制器路由执行器
type FunctionRouterExecutor struct {
	router *FunctionRouter //函数路由
}

func (this *FunctionRouterExecutor) Exec(context RouterContext) {
	this.router.function(context)
}
