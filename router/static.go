package router

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// 静态文件路由
// 静态文件方法路由仅用于隔离路由空间,本身并不具备任何功能
type StaticRouter struct {
	SpaceRouter
	path string //当前静态路由对应的本地文件目录
}

// Pass 传递指定的路由环境给当前的路由器
// context: 上下文环境
// return: 返回路由是否处理了该请求
// 如果请求已经被处理了,则该请求不应该继续被传递
func (this *StaticRouter) Pass(context RouterContext) bool {
	//只响应GET请求
	if strings.EqualFold(context.Method(), "GET") {
		var parts = context.RouterParts()
		var currentPath = ""
		for index := this.Level(); index < len(parts); index++ {
			var param = parts[index]
			if param != "" && param != ".." {
				filepath.Join(currentPath, param)
			}
		}
		currentPath = filepath.Join(this.path, currentPath)
		var info, err = os.Stat(currentPath)
		if err == nil && !info.IsDir() {
			var sre = &StaticRouterExecutor{currentPath}
			context.AddContextExecutor(sre)
			return true
		}
	}
	return false
}

// 静态文件执行器
type StaticRouterExecutor struct {
	path string //文件路径
}

func (this *StaticRouterExecutor) Exec(context RouterContext) {
	http.ServeFile(context.ResponseWriter(), context.Request(), this.path)
}
