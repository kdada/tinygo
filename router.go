package tinygo

import (
	"path/filepath"

	"github.com/kdada/tinygo/router"
)

// 根路由,其他路由应该作为根路由的子路由
var RootRouter = router.NewRootRouter()

// 生成静态路由
func generateStaticRouters() {
	//生成favicon.ico跳转路由
	RootRouter.AddChild(router.NewStaticRouter("favicon.ico", "favicon.ico"))
	//其他静态路由
	for _, path := range tinyConfig.static {
		var components = filepath.SplitList(path)
		var count = len(components)
		var lastRouter = RootRouter
		if count > 2 {
			for _, c := range components {
				var space = router.NewSpaceRouter(c)
				lastRouter.AddChild(space)
				lastRouter = space
			}
		}
		lastRouter.AddChild(router.NewStaticRouter(components[count-1], path))
	}
}

// 控制器方法路由信息
type RouterInfo struct {
	methodName string     //方法名
	httpMethod HttpMethod //http方法
	extensions []string   //url扩展
}

// Info 获取路由信息
func (this *RouterInfo) Info() (string, string, []string) {
	return this.methodName, string(this.httpMethod), this.extensions
}

// NewRouterInfo 创建单个路由信息
func NewRouterInfo(name string, method HttpMethod, extensions []string) interface{} {
	return &RouterInfo{name, method, extensions}
}

// SetHomePage 设置首页
func SetHomePage(path string) bool {
	var root, ok = RootRouter.(*router.SpaceRouter)
	if ok {
		root.SetDefaultPage(path)
	}
	return ok
}

// NewSpaceRouter 创建空间路由
//  name:路由名称
func NewSpaceRouter(name string) router.Router {
	return router.NewSpaceRouter(name)
}

// NewStaticRouter 创建控制器方法路由
//  name:静态路由名
//  path:静态文件本地目录
// 例如
//  name "static"
//  path "content/static/"
// 即url static/css/index.css 映射为本地目录 content/static/css/index.css
func NewStaticRouter(name, path string) router.Router {
	return router.NewStaticRouter(name, path)
}

// NewControllerRouter 创建控制器方法路由
//  instance:结构体实例,必须是结构体指针,并且在Routers方法中返回方法路由信息
func NewControllerRouter(instance router.Controller) router.Router {
	return router.NewControllerRouter(instance)
}

// NewRestfulControllerRouter 创建Restful控制器路由
//  instance:结构体实例,必须是结构体指针
func NewRestfulControllerRouter(instance router.RestfulController) router.Router {
	return router.NewRestfulControllerRouter(instance)
}

// NewFunctionRouter 创建函数路由
func NewFunctionRouter(name string, f func(router.RouterContext)) router.Router {
	return router.NewFunctionRouter(name, f)
}
