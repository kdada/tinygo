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
	BaseRouter
	path string //当前静态路由对应的本地文件目录
}

// Pass 传递指定的路由环境给当前的路由器
//  context: 上下文环境
//  return: 返回路由是否处理了该请求
// 如果请求已经被处理了,则该请求不应该继续被传递
func (this *StaticRouter) Pass(context RouterContext) bool {
	//只响应GET请求
	if strings.EqualFold(context.Method(), "GET") {
		var sre = &StaticRouterExecutor{this.Level(), this.path}
		context.AddRouter(this)
		context.AddContextExecutor(sre)
		//设置是静态路由
		context.SetStatic(true)
		return true
	}
	return false
}

// 静态文件执行器
type StaticRouterExecutor struct {
	level int    //路由层级
	path  string //文件路径
}

func (this *StaticRouterExecutor) Exec(context RouterContext) {
	var parts = context.RouterParts()
	var currentPath = ""
	var malice = false
	for index := this.level + 1; index < len(parts); index++ {
		var param = parts[index]
		if strings.Contains(param, `\`) {
			//param不应该包含这种分隔符
			malice = true
			break
		}
		if param != "" && param != ".." {
			currentPath = filepath.Join(currentPath, param)
		}
	}
	if !malice {
		currentPath = filepath.Join(this.path, currentPath)
		var info, err = os.Stat(currentPath)
		if err == nil && !info.IsDir() {
			http.ServeFile(context.ResponseWriter(), context.Request(), currentPath)
			return
		}
	}
	//空文件路径
	http.NotFound(context.ResponseWriter(), context.Request())
}
