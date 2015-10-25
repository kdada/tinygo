package router

import "strings"

// 空间路由
// 空间路由仅用于隔离路由空间,本身并不具备任何功能
type SpaceRouter struct {
	BaseRouter
	defaultPage string //默认页面
}

// 默认页面,当路由段到当前空间路由为止的时候
// 当前路由将使用默认页面继续查找路由
func (this *SpaceRouter) DefaultPage() string {
	return this.defaultPage
}

// 设置默认页面
func (this *SpaceRouter) SetDefaultPage(ref string) {
	this.defaultPage = ref
}

// Pass 传递指定的路由环境给当前的路由器
//  context: 上下文环境
//  return: 返回路由是否处理了该请求
// 如果请求已经被处理了,则该请求不应该继续被传递
func (this *SpaceRouter) Pass(context RouterContext) bool {
	var parts = context.RouterParts()
	if len(parts) == (this.level+1) && this.defaultPage != "" {
		//默认路由分段后添加默认分段
		var ps = strings.Split(strings.Trim(this.defaultPage, "/"), "/")
		if len(ps) > 0 && ps[0] != "" {
			context.SetRouterParts(append(parts, ps...))
			return this.Pass(context)
		}
		return false
	}
	if len(parts) > (this.level + 1) {
		//当前路由段名称
		var name = parts[this.level]
		//检查当前路由能否处理
		var routeData, canExec = this.check(name)
		if canExec {
			//子路由段名称
			var pathName = parts[this.level+1]
			var childRouter, ok = this.children[pathName]
			if ok {
				//直接传递
				ok = childRouter.Pass(context)
			} else {
				//遍历正则路由
				for _, v := range this.regchildren {
					ok = v.Pass(context)
					if ok {
						break
					}
				}
			}
			if ok {
				//添加路由参数
				for k, v := range routeData {
					context.AddRouterParams(k, v)
				}
				//添加当前路由
				context.AddRouter(this)
			}
			return ok
		}
	}
	return false
}
