package tinygo

import (
	"path/filepath"

	"github.com/kdada/tinygo/router"
)

// 路由信息类型
type RouterInfos []interface{}

// 路由参数段类型
type RouterParams []string

// 根路由,其他路由应该作为根路由的子路由
var rootRouter = router.NewRootRouter()

// AddRouter 添加router到跟路由
func AddRouter(router router.Router) {
	rootRouter.AddChild(router)
}

// 生成静态路由
func generateStaticRouters() {
	//生成favicon.ico跳转路由
	AddRouter(router.NewStaticRouter("favicon.ico", "favicon.ico"))
	//其他静态路由
	for _, path := range tinyConfig.static {
		var components = filepath.SplitList(path)
		var count = len(components)
		var lastRouter = rootRouter
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
	aliasName  string     //方法别名(可包含正则)
	httpMethod HttpMethod //http方法
	extensions []string   //url扩展
}

// Info 获取路由信息
func (this *RouterInfo) Info() (string, string, string, []string) {
	if this.aliasName == "" {
		this.aliasName = this.methodName
	}
	return this.methodName, this.aliasName, string(this.httpMethod), this.extensions
}

// NewRouterInfo 创建单个路由信息
//  name:方法名,区分大小写
//  method:http方法
//  extensions:扩展内容,包含正则
func NewRouterInfo(name string, method HttpMethod, extensions []string) interface{} {
	return &RouterInfo{name, name, method, extensions}
}

// NewAliasRouterInfo 创建单个路由信息
//  name:方法名,区分大小写
//  alias:方法别名,可包含正则,例如:info_(id=\d+).html
//  method:http方法
//  extensions:扩展内容,包含正则
func NewAliasRouterInfo(name, alias string, method HttpMethod, extensions []string) interface{} {
	return &RouterInfo{name, alias, method, extensions}
}

// NewGetRouterInfo 创建单个GET路由信息
//  name:方法名,区分大小写
//  extensions:扩展内容,包含正则
func NewGetRouterInfo(name string, extensions []string) interface{} {
	return &RouterInfo{name, name, HttpMethodGet, extensions}
}

// NewAliasGetRouterInfo 创建单个GET路由信息
//  name:方法名,区分大小写
//  alias:方法别名,可包含正则,例如:info_(id=\d+).html
//  extensions:扩展内容,包含正则
func NewAliasGetRouterInfo(name, alias string, extensions []string) interface{} {
	return &RouterInfo{name, alias, HttpMethodGet, extensions}
}

// NewPostRouterInfo 创建单个POST路由信息
//  name:方法名,区分大小写
//  extensions:扩展内容,包含正则
func NewPostRouterInfo(name string, extensions []string) interface{} {
	return &RouterInfo{name, name, HttpMethodPost, extensions}
}

// NewAliasPostRouterInfo 创建单个POST路由信息
//  name:方法名,区分大小写
//  alias:方法别名,可包含正则,例如:info_(id=\d+).html
//  extensions:扩展内容,包含正则
func NewAliasPostRouterInfo(name, alias string, extensions []string) interface{} {
	return &RouterInfo{name, alias, HttpMethodPost, extensions}
}

// SetHomePage 设置首页
func SetHomePage(path string) bool {
	var root, ok = rootRouter.(*router.SpaceRouter)
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

// NewComplexControllerRouter 创建使用指定路径以及路由信息的控制器路由
//  path:url路径,例如/list/home,支持url正则
//  instance:结构体实例,必须是结构体指针
//  info:路由信息
func NewComplexControllerRouter(path string, instance router.Controller, info []interface{}) router.Router {
	return router.NewComplexControllerRouter(path, instance, info)
}

// NewRestfulControllerRouter 创建Restful控制器路由
//  instance:结构体实例,必须是结构体指针
func NewRestfulControllerRouter(instance router.RestfulController) router.Router {
	return router.NewRestfulControllerRouter(instance)
}

// NewPathRestfulControllerRouter 将instance作为指定path的处理控制器
//  path:url路径,例如/list/(id).html,支持url正则
//  instance:结构体实例,必须是结构体指针
func NewPathRestfulControllerRouter(path string, instance router.RestfulController) router.Router {
	return router.NewPathRestfulControllerRouter(path, instance)
}

// NewFunctionRouter 创建函数路由
//  path:url路径,例如/list/(id).html,支持url正则
//  f:执行方法
func NewFunctionRouter(path string, f func(router.RouterContext)) router.Router {
	return router.NewFunctionRouter(path, f)
}
