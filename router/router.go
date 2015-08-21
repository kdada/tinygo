package router

import (
	"strings"
)

//路由器接口
type IRouter interface {
	// Name 返回当前路由的名称
	Name() string
	// Super 返回当前路由的父路由
	Super() IRouter
	// SetSuper 设置父路由
	SetSuper(super IRouter)
	// Level 返回当前路由层级
	Level() int
	// SetLevel 设置当前路由层级
	SetLevel(level int)
	// Check 检查当前路由器是否支持该请求
	// routes:地址段
	// 例如:请求的url:/admin/getsomething/1234?other=xxx
	// routes: ["admin","getsomething","1234"]
	// return:如果为true表示当前路由器可以处理该路由
	//Check(routes []string) bool
	// Pass 传递指定的路由环境给当前的路由器
	// context: 上下文环境
	// return: 返回路由是否处理了该请求
	// 如果请求已经被处理了,则该请求不应该继续被传递
	Pass(context IRouterContext) bool
	// Child 通过名称获取子路由
	Child(name string) (IRouter, bool)
	// AddChild 添加子路由
	AddChild(router IRouter) bool
	// AddChildren 批量添加添加子路由,如果已经存在同名路由,则添加失败
	AddChildren(routers ...IRouter) bool
	// RemoveChild 移除子路由
	RemoveChild(name string) bool
	// AddBeforeFilter 添加前置过滤器
	AddBeforeFilter(filter IRouterFilter) bool
	// RemoveBeforeFilter 移除后置过滤器
	RemoveBeforeFilter(filter IRouterFilter) bool
	// AddAfterFilter 添加后置过滤器
	AddAfterFilter(fileter IRouterFilter) bool
	// RemoveAfterFilter 移除后置过滤器
	RemoveAfterFilter(filter IRouterFilter) bool
}

//基础路由器数据
type BaseRouter struct {
	name          string             //当前路由名称
	super         IRouter            //上级路由
	level         int                //路由层级
	children      map[string]IRouter //子路由
	beforeFilters []IRouterFilter    //在子路由处理之前执行的过滤器
	afterFilters  []IRouterFilter    //在子路由处理之后执行的过滤器
}

// Init 初始化基础路由数据
func (this *BaseRouter) Init(name string) {
	this.name = name
	this.children = make(map[string]IRouter, 0)
	this.beforeFilters = make([]IRouterFilter, 0)
	this.afterFilters = make([]IRouterFilter, 0)
}

// Name 返回当前路由的名称
func (this *BaseRouter) Name() string {
	return this.name
}

// Super 返回当前路由的父路由
func (this *BaseRouter) Super() IRouter {
	return this.super
}

// SetSuper 设置父路由
func (this *BaseRouter) SetSuper(super IRouter) {
	this.super = super
}

// Level 返回当前路由层级
func (this *BaseRouter) Level() int {
	return this.level
}

// SetLevel 设置当前路由层级
func (this *BaseRouter) SetLevel(level int) {
	this.level = level
}

// Check 检查当前路由器是否支持该请求
// routes:地址段
// 请求的url:/admin/getsomething/1234?other=xxx
// route: ["admin","getsomething","1234"]
// return:如果为true表示当前路由器可以处理该路由
//func (this *BaseRouter) Check(routes []string) bool {
//	if len(routes) > this.level {
//		return strings.ToLower(routes[this.level]) == strings.ToLower(this.name)
//	}
//	return false
//}

// Pass 传递指定的路由环境给当前的路由器
// context: 上下文环境
// return: 返回路由是否处理了该请求
// 如果请求已经被处理了,则该请求不应该继续被传递
func (this *BaseRouter) Pass(context IRouterContext) bool {
	return false
}

// Child 通过名称获取子路由
func (this *BaseRouter) Child(name string) (IRouter, bool) {
	var router, ok = this.children[strings.ToLower(name)]
	if ok {
		return router, true
	}
	return nil, false
}

// AddChild 添加子路由,如果已经存在同名路由,则添加失败
func (this *BaseRouter) AddChild(router IRouter) bool {
	if router != nil {
		var _, ok = this.children[strings.ToLower(router.Name())]
		if !ok {
			this.children[strings.ToLower(router.Name())] = router
			router.SetSuper(this)
			router.SetLevel(this.level + 1)
			return true
		}
	}
	return false
}

// AddChildren 批量添加添加子路由,如果已经存在同名路由,则添加失败
func (this *BaseRouter) AddChildren(routers ...IRouter) bool {
	for _, router := range routers {
		this.AddChild(router)
	}
	return false
}

// RemoveChild 移除子路由
// name:子路由名称
func (this *BaseRouter) RemoveChild(name string) bool {
	var _, ok = this.children[name]
	if ok {
		delete(this.children, name)
		return true
	}
	return false
}

// AddBeforeFilter 添加前置过滤器
func (this *BaseRouter) AddBeforeFilter(filter IRouterFilter) bool {
	if filter != nil {
		this.beforeFilters = append(this.beforeFilters, filter)
		return true
	}
	return false
}

// RemoveBeforeFilter 移除前置过滤器
func (this *BaseRouter) RemoveBeforeFilter(filter IRouterFilter) bool {
	for index, child := range this.beforeFilters {
		if child == filter {
			this.beforeFilters = append(this.beforeFilters[:index], this.beforeFilters[index+1:]...)
			return true
		}
	}
	return false
}

// ExecBeforeFilter 过滤请求
// return:返回true表示继续处理,否则终止路由过程
func (this *BaseRouter) ExecBeforeFilter(context IRouterContext) bool {
	for _, router := range this.beforeFilters {
		var goon = router.Filter(context)
		if !goon {
			return false
		}
	}
	return true
}

// AddAfterFilter 添加后置过滤器
func (this *BaseRouter) AddAfterFilter(filter IRouterFilter) bool {
	if filter != nil {
		this.afterFilters = append(this.afterFilters, filter)
		return true
	}
	return false
}

// RemoveAfterFilter 移除后置过滤器
func (this *BaseRouter) RemoveAfterFilter(filter IRouterFilter) bool {
	for index, child := range this.afterFilters {
		if child == filter {
			this.afterFilters = append(this.afterFilters[:index], this.afterFilters[index+1:]...)
			return true
		}
	}
	return false
}

// ExecAfterFilter 过滤请求
// return:返回true表示继续处理,否则终止路由过程
func (this *BaseRouter) ExecAfterFilter(context IRouterContext) bool {
	for _, router := range this.afterFilters {
		var goon = router.Filter(context)
		if !goon {
			return false
		}
	}
	return true
}
