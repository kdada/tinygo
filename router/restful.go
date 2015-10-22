package router

import (
	"reflect"
	"strings"
)

// Restful控制器路由
// Restful控制器路由仅支持实现了RestfulController接口的控制器
type RestfulRouter struct {
	BaseRouter
	instanceType reflect.Type //结构体类型
}

// Pass 传递指定的路由环境给当前的路由器
// context: 上下文环境
// return: 返回路由是否处理了该请求
// 如果请求已经被处理了,则该请求不应该继续被传递
func (this *RestfulRouter) Pass(context RouterContext) bool {
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
			var rre = &RestfulRouterExecutor{this}
			context.AddRouter(this)
			context.AddContextExecutor(rre)
			return true
		}
	}

	return false
}

// Restful控制器路由执行器
type RestfulRouterExecutor struct {
	router *RestfulRouter //Restful控制器路由
}

func (this *RestfulRouterExecutor) Exec(context RouterContext) {
	//创建类型实例并执行相应方法
	//使用method进行方法调用
	var instance = reflect.New(this.router.instanceType)
	controller, ok := instance.Interface().(RestfulController)
	if ok {
		controller.SetContext(context)
		controller.SetRouter(this.router)
		var method = strings.ToUpper(context.Method())
		switch method {
		case "GET":
			controller.Get()
		case "POST":
			controller.Post()
		case "PUT":
			controller.Put()
		case "DELETE":
			controller.Delete()
		}
	}
}
