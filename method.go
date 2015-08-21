package tinygo

import (
	"reflect"
	"strings"
	"github.com/kdada/tinygo/info"
	"github.com/kdada/tinygo/router"
	"github.com/kdada/tinygo/session"
)

// 控制器方法路由信息
type RouterInfo struct {
	MethodName string          //方法名
	HttpMethod info.HttpMethod //http方法
	Extensions []string        //url扩展
}

// 控制器方法路由
// 控制器方法路由仅支持实现了IController接口的控制器
type MethodRouter struct {
	router.SpaceRouter
	instanceType reflect.Type    //结构体类型
	methodType   reflect.Method  //方法类型
	httpMethod   info.HttpMethod //能处理的Http方法
	extensions   []string        //路由扩展
}

// NewControllerRouter 创建控制器方法路由
// instance:结构体实例,必须是结构体指针,并且在Routers方法中返回方法路由信息
func NewControllerRouter(instance IController) router.IRouter {
	var ptrType = reflect.TypeOf(instance)
	if ptrType.Kind() == reflect.Ptr && ptrType.Elem().Kind() == reflect.Struct {
		var instanceType = ptrType.Elem()
		//路由信息
		var routersInfo = instance.Routers()
		//控制器名
		var spaceName = instanceType.Name()
		spaceName = strings.TrimSuffix(spaceName, "Controller")
		var controllerRouter = router.NewSpaceRouter(spaceName)
		for _, routerInfo := range routersInfo {
			var methodRouter = new(MethodRouter)
			methodRouter.Init(routerInfo.MethodName)
			methodRouter.httpMethod = routerInfo.HttpMethod
			methodRouter.extensions = routerInfo.Extensions
			methodRouter.instanceType = instanceType
			var ok = false
			methodRouter.methodType, ok = ptrType.MethodByName(routerInfo.MethodName)
			if ok {
				controllerRouter.AddChild(methodRouter)
			}
		}
		return controllerRouter
	}
	return nil
}

// Pass 传递指定的路由环境给当前的路由器
// context: 上下文环境
// return: 返回路由是否处理了该请求
// 如果请求已经被处理了,则该请求不应该继续被传递
func (this *MethodRouter) Pass(context router.IRouterContext) bool {
	var httpContext, ok = context.(*router.HttpContext)
	if ok {
		if string(this.httpMethod) == httpContext.Request.Method {
			var parts = context.RouterParts()
			if len(parts) > this.Level()+len(this.extensions) {
				//将路由中多余的部分作为查询参数添加到http环境中
				for index := this.Level() + 1; index < len(parts); index++ {
					httpContext.AddParams(this.extensions[index], parts[index])
				}
			}
			//添加Session信息
			var cookie, err = httpContext.Request.Cookie(info.DefaultSessionCookieName)
			var ss session.ISession
			var ok bool = false
			if err == nil {
				ss, ok = SessionProvider.Session(cookie.Value)
			}
			if !ok {
				ss, ok = SessionProvider.CreateSession()
			}
			if ok {
				httpContext.Session = ss
				httpContext.ResponseWriter.Header().Set("Set-Cookie", info.DefaultSessionCookieName+"="+ss.SessionId())
			}
			//创建类型实例并执行相应方法
			//使用method进行方法调用
			var instance = reflect.New(this.instanceType)
			controller, ok := instance.Interface().(IController)
			if ok {
				controller.SetContext(httpContext)
				controller.SetRouter(this)
				var value = []reflect.Value{instance}
				//执行前置过滤器
				ok = this.ExecBeforeFilter(context)
				if ok {
					//执行处理方法
					this.methodType.Func.Call(value)
					//执行后置过滤器
					ok = this.ExecAfterFilter(context)
					return ok
				}
			}
		}
	}
	return false
}
