// Package router 实现了tinygo的路由系统
package router

import (
	"net/http"
)

// 路由器接口
type Router interface {
	// Name 返回当前路由的名称
	Name() string
	// Super 返回当前路由的父路由
	Super() Router
	// SetSuper 设置父路由
	SetSuper(super Router)
	// Level 返回当前路由层级
	Level() int
	// Reg 返回当前路由是否是正则路由
	Reg() bool
	// SetLevel 设置当前路由层级
	SetLevel(level int)
	// Pass 传递指定的路由环境给当前的路由器
	Pass(context RouterContext) bool
	// Child 通过名称获取子路由
	Child(name string) (Router, bool)
	// AddChild 添加子路由
	AddChild(router Router) bool
	// AddChildren 批量添加添加子路由,如果已经存在同名路由,则添加失败
	AddChildren(routers ...Router) bool
	// RemoveChild 移除子路由
	RemoveChild(name string) bool
	// AddBeforeFilter 添加前置过滤器
	AddBeforeFilter(filter RouterFilter) bool
	// RemoveBeforeFilter 移除前置过滤器
	RemoveBeforeFilter(filter RouterFilter) bool
	// ExecBeforeFilter 执行前置过滤器
	ExecBeforeFilter(context RouterContext) bool
	// AddAfterFilter 添加后置过滤器
	AddAfterFilter(fileter RouterFilter) bool
	// RemoveAfterFilter 移除后置过滤器
	RemoveAfterFilter(filter RouterFilter) bool
	// ExecAfterFilter 执行后置过滤器
	ExecAfterFilter(context RouterContext) bool
}

// 默认页面路由器接口
type DefaultPageRouter interface {
	Router
	// 默认页面
	DefaultPage() string
	// 设置默认页面
	SetDefaultPage(ref string)
}

// 路由过滤器
type RouterFilter interface {
	// Filter 过滤该请求
	// return:返回true表示继续处理,否则终止路由过程,后续的过滤器也不会执行
	Filter(context RouterContext) bool
}

// Context执行器
type ContextExecutor interface {
	Exec(context RouterContext)
}

// 路由环境接口
type RouterContext interface {
	// Method 返回请求的HTTP方法
	Method() string
	// responseWriter 返回responseWriter
	ResponseWriter() http.ResponseWriter
	// Request 返回Request
	Request() *http.Request
	// RouterParts 返回路由段
	RouterParts() []string
	// SetRouterParts 设置路由段
	SetRouterParts(parts []string)
	// Static 返回是否是静态路由
	Static() bool
	// SetStatic 设置当前上下文为静态路由上下文
	SetStatic(static bool)
	// AddRouterParams 添加路由参数
	AddRouterParams(key, value string)
	// RemoveRouterParams 移除路由参数
	RemoveRouterParams(key string)
	// AddRouter 添加执行路由,最后一级路由最先添加
	AddRouter(router Router)
	// AddContextExector 添加执行器
	AddContextExecutor(exector ContextExecutor)
}
