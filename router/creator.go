package router

import (
	"fmt"
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
//  name:路由名称
func NewSpaceRouter(name string) Router {
	var router = new(SpaceRouter)
	router.Init(name)
	return router
}

// NewStaticRouter 创建控制器方法路由
//  name:静态路由名
//  path:静态文件本地目录
// 例如
//  name "static"
//  path "content/static/"
// 即url static/css/index.css 映射为本地目录 content/static/css/index.css
func NewStaticRouter(name, path string) Router {
	var router = new(StaticRouter)
	router.Init(name)
	router.path = strings.TrimRight(path, "/")
	return router
}

// NewControllerRouter 创建控制器方法路由
//  instance:结构体实例,必须是结构体指针,并且在Routers方法中返回方法路由信息
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
				var name, alias, method, extensions = info.Info()
				var regs, err = ParseRegs(extensions)
				if err == nil {
					var methodRouter = new(MethodRouter)
					methodRouter.Init(alias)
					methodRouter.httpMethod = method
					methodRouter.extensions = regs
					methodRouter.instanceType = instanceType
					var ok = false
					methodRouter.methodType, ok = ptrType.MethodByName(name)
					if ok {
						controllerRouter.AddChild(methodRouter)
					}
				} else {
					fmt.Println(err)
					return nil
				}

			}
		}
		return controllerRouter
	}
	return nil
}

// NewComplexControllerRouter 创建使用指定路径以及路由信息的控制器路由
//  path:url路径,例如/list/home,支持url正则
//  instance:结构体实例,必须是结构体指针
//  info:路由信息
func NewComplexControllerRouter(path string, instance Controller, info []interface{}) Router {
	var ptrType = reflect.TypeOf(instance)
	if ptrType.Kind() == reflect.Ptr && ptrType.Elem().Kind() == reflect.Struct {
		var instanceType = ptrType.Elem()
		//控制器名
		var spaces = strings.Split(strings.Trim(path, "/"), "/")
		var spacesLen = len(spaces)
		if spacesLen > 0 {
			var tail = NewSpaceRouter(spaces[spacesLen-1]) //最后一级路由,即控制器路由
			var header = tail                              //一级路由
			for i := spacesLen - 2; i >= 0; i-- {
				var superRouter = NewSpaceRouter(spaces[i])
				superRouter.AddChild(header)
				header = superRouter
			}
			//根据路由信息创建方法路由
			for _, routerInfo := range info {
				var ri, ok = routerInfo.(ControllerRouter)
				if ok {
					var name, alias, method, extensions = ri.Info()
					var regs, err = ParseRegs(extensions)
					if err == nil {
						var methodRouter = new(MethodRouter)
						methodRouter.Init(alias)
						methodRouter.httpMethod = method
						methodRouter.extensions = regs
						methodRouter.instanceType = instanceType
						var ok = false
						methodRouter.methodType, ok = ptrType.MethodByName(name)
						if ok {
							tail.AddChild(methodRouter)
						}
					} else {
						fmt.Println(err)
						return nil
					}
				} else {
					return nil
				}
			}
			return header
		}
	}
	return nil
}

// NewRestfulControllerRouter 创建Restful控制器路由
//  instance:结构体实例,必须是结构体指针
func NewRestfulControllerRouter(instance RestfulController) Router {
	var ptrType = reflect.TypeOf(instance)
	if ptrType.Kind() == reflect.Ptr && ptrType.Elem().Kind() == reflect.Struct {
		var instanceType = ptrType.Elem()
		//控制器名
		var spaceName = instanceType.Name()
		spaceName = strings.TrimSuffix(spaceName, "Controller")
		var controllerRouter = new(RestfulRouter)
		controllerRouter.Init(spaceName)
		controllerRouter.instanceType = instanceType
		return controllerRouter
	}
	return nil
}

// NewPathRestfulControllerRouter 将instance作为指定path的处理控制器
//  path:url路径,例如/list/(id).html,支持url正则
//  instance:结构体实例,必须是结构体指针
func NewPathRestfulControllerRouter(path string, instance RestfulController) Router {
	var spaces = strings.Split(strings.Trim(path, "/"), "/")
	var spacesLen = len(spaces)
	if spacesLen > 0 {
		var tail = NewRestfulControllerRouter(instance)
		var tailRouter = tail.(*RestfulRouter)
		tailRouter.setName(spaces[spacesLen-1])
		for i := spacesLen - 2; i >= 0; i-- {
			var superRouter = NewSpaceRouter(spaces[i])
			superRouter.AddChild(tail)
			tail = superRouter
		}
		return tail
	}
	return nil
}

// NewFunctionRouter 创建函数路由
func NewFunctionRouter(path string, f func(RouterContext)) Router {
	var spaces = strings.Split(strings.Trim(path, "/"), "/")
	var spacesLen = len(spaces)
	if spacesLen > 0 && f != nil {
		var tail = new(FunctionRouter) //最后一级路由,即函数路由
		tail.Init(spaces[spacesLen-1])
		tail.function = f
		var header Router = tail //一级路由
		for i := spacesLen - 2; i >= 0; i-- {
			var superRouter = NewSpaceRouter(spaces[i])
			superRouter.AddChild(header)
			header = superRouter
		}
		return header
	}
	return nil
}
