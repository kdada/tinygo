package router

import (
	"reflect"
	"strings"
)

// 控制器方法路由
// 控制器方法路由仅支持实现了Controller接口的控制器
type MethodRouter struct {
	BaseRouter
	instanceType reflect.Type   //结构体类型
	methodType   reflect.Method //方法类型
	httpMethod   string         //能处理的Http方法
	extensions   []string       //路由扩展
}

// Pass 传递指定的路由环境给当前的路由器
// context: 上下文环境
// return: 返回路由是否处理了该请求
// 如果请求已经被处理了,则该请求不应该继续被传递
func (this *MethodRouter) Pass(context RouterContext) bool {
	if strings.EqualFold(string(this.httpMethod), context.Method()) {
		var parts = context.RouterParts()
		if len(parts) == this.Level()+len(this.extensions)+1 {
			//检查当前路由是否能够处理
			var route = parts[this.Level()]
			var routeData, ok = this.check(route)
			if ok {
				//将路由中多余的部分作为查询参数添加到http环境中
				for index := this.Level() + 1; index < len(parts); index++ {
					context.AddRouterParams(this.extensions[index], parts[index])
				}
				//添加路由参数
				for k, v := range routeData {
					context.AddRouterParams(k, v)
				}
				//处理
				var mre = &MethodRouterExecutor{this}
				context.AddRouter(this)
				context.AddContextExecutor(mre)
				return true
			}
		}
	}
	return false
}

// 方法路由执行器
type MethodRouterExecutor struct {
	router *MethodRouter //方法路由
}

func (this *MethodRouterExecutor) Exec(context RouterContext) {
	//创建类型实例并执行相应方法
	//使用method进行方法调用
	var instance = reflect.New(this.router.instanceType)
	controller, ok := instance.Interface().(Controller)
	if ok {
		controller.SetContext(context)
		controller.SetRouter(this.router)
		var value = []reflect.Value{instance}
		//执行处理方法
		this.router.methodType.Func.Call(value)
	}
}
