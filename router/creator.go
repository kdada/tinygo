package router

import (
	"reflect"
	"strings"
)

// NewRootRouter 创建根路由
func NewRootRouter() Router {
	var router = new(SpaceRouter)
	router.Init("")
	return router
}

// NewSpaceRouter 创建空间路由
// name:路由名称
func NewSpaceRouter(name string) Router {
	var router = new(SpaceRouter)
	router.Init(name)
	return router
}

// NewStaticRouter 创建控制器方法路由
// name:静态路由名
// path:静态文件本地目录
// 例如
// name "static"
// path "content/static/"
// 即url static/css/index.css 映射为本地目录 content/static/css/index.css
func NewStaticRouter(name, path string) Router {
	var router = new(StaticRouter)
	router.Init(name)
	router.path = strings.TrimRight(path, "/")
	return router
}

// NewControllerRouter 创建控制器方法路由
// instance:结构体实例,必须是结构体指针,并且在Routers方法中返回方法路由信息
func NewControllerRouter(instance Controller) Router {
	var ptrType = reflect.TypeOf(instance)
	if ptrType.Kind() == reflect.Ptr && ptrType.Elem().Kind() == reflect.Struct {
		var instanceType = ptrType.Elem()
		//路由信息
		var routersInfo = instance.Routers()
		//控制器名
		var spaceName = instanceType.Name()
		spaceName = strings.TrimSuffix(spaceName, "Controller")
		var controllerRouter = NewSpaceRouter(spaceName)
		for _, routerInfo := range routersInfo {
			var info, ok = routerInfo.(ControllerRouter)
			if ok {
				var name, method, extensions = info.Info()
				var methodRouter = new(MethodRouter)
				methodRouter.Init(name)
				methodRouter.httpMethod = method
				methodRouter.extensions = extensions
				methodRouter.instanceType = instanceType
				var ok = false
				methodRouter.methodType, ok = ptrType.MethodByName(name)
				if ok {
					controllerRouter.AddChild(methodRouter)
				}
			}
		}
		return controllerRouter
	}
	return nil
}

// NewRestfulControllerRouter 创建Restful控制器路由
// instance:结构体实例,必须是结构体指针
func NewRestfulControllerRouter(instance RestfulController) Router {
	var ptrType = reflect.TypeOf(instance)
	if ptrType.Kind() == reflect.Ptr && ptrType.Elem().Kind() == reflect.Struct {
		var instanceType = ptrType.Elem()
		//控制器名
		var spaceName = instanceType.Name()
		spaceName = strings.TrimSuffix(spaceName, "Controller")
		var controllerRouter = new(RestfulRouter)
		controllerRouter.instanceType = instanceType
		return controllerRouter
	}
	return nil
}

// NewFunctionRouter 创建函数路由
func NewFunctionRouter(name string, f func(RouterContext)) Router {
	var router = new(FunctionRouter)
	router.Init(name)
	router.function = f
	return router
}
