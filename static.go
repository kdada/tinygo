package tinygo

import (
	"net/http"
	"os"
	"strings"

	"github.com/kdada/tinygo/info"
	"github.com/kdada/tinygo/router"
)

// 静态文件路由
// 静态文件方法路由仅用于隔离路由空间,本身并不具备任何功能
type StaticRouter struct {
	router.SpaceRouter
	path string //当前静态路由对应的本地文件目录
}

// NewStaticRouter 创建控制器方法路由
// name:静态路由名
// path:静态文件本地目录
// 例如
// name "static"
// path "content/static/"
// 即url static/css/index.css 映射为本地目录 content/static/css/index.css
func NewStaticRouter(name, path string) router.Router {
	var router = new(StaticRouter)
	router.Init(name)
	router.path = strings.TrimRight(path, "/")
	return router
}

// Pass 传递指定的路由环境给当前的路由器
// context: 上下文环境
// return: 返回路由是否处理了该请求
// 如果请求已经被处理了,则该请求不应该继续被传递
func (this *StaticRouter) Pass(context router.RouterContext) bool {
	var httpContext, ok = context.(*router.HttpContext)
	//只响应GET请求
	if ok && httpContext.Request.Method == string(info.HttpMethodGet) {
		var currentPath = ""
		for index := this.Level(); index < len(httpContext.UrlParts); index++ {
			var param = httpContext.UrlParts[index]
			if param != "" && param != ".." {
				currentPath += "/" + param
			}
		}
		currentPath = this.path + currentPath
		var info, err = os.Stat(currentPath)
		if err == nil && !info.IsDir() {
			http.ServeFile(httpContext.ResponseWriter, httpContext.Request, currentPath)
			return true
		}
	}
	HttpNotFound(httpContext.ResponseWriter, httpContext.Request)
	return true
}
