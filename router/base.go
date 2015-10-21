package router

import "strings"

//基础路由器数据
type BaseRouter struct {
	name          string            //当前路由名称
	super         Router            //上级路由
	level         int               //路由层级
	segment       *RegSegment       //正则路由信息
	reg           bool              //是否是正则路由
	children      map[string]Router //子路由
	regchildren   map[string]Router //正则子路由
	beforeFilters []RouterFilter    //在子路由处理之前执行的过滤器
	afterFilters  []RouterFilter    //在子路由处理之后执行的过滤器
}

// Init 初始化基础路由数据
func (this *BaseRouter) Init(name string) {
	var segment, err = ParseReg(name)
	if err == nil {
		this.segment = segment
		this.reg = true
	}
	this.name = name
	this.children = make(map[string]Router, 0)
	this.regchildren = make(map[string]Router, 0)
	this.beforeFilters = make([]RouterFilter, 0)
	this.afterFilters = make([]RouterFilter, 0)
}

// Name 返回当前路由的名称
func (this *BaseRouter) Name() string {
	return this.name
}

// Super 返回当前路由的父路由
func (this *BaseRouter) Super() Router {
	return this.super
}

// SetSuper 设置父路由
func (this *BaseRouter) SetSuper(super Router) {
	this.super = super
}

// Level 返回当前路由层级
func (this *BaseRouter) Level() int {
	return this.level
}

// Reg 返回当前路由是否是正则路由
func (this *BaseRouter) Reg() bool {
	return this.reg
}

// SetLevel 设置当前路由层级
func (this *BaseRouter) SetLevel(level int) {
	this.level = level
	for _, child := range this.children {
		child.SetLevel(this.level + 1)
	}
	for _, child := range this.regchildren {
		child.SetLevel(this.level + 1)
	}
}

// Pass 传递指定的路由环境给当前的路由器
// context: 上下文环境
// return: 返回路由是否处理了该请求
// 如果请求已经被处理了,则该请求不应该继续被传递
func (this *BaseRouter) Pass(context RouterContext) bool {
	return false
}

// Child 通过名称获取子路由
func (this *BaseRouter) Child(name string) (Router, bool) {
	name = strings.ToLower(name)
	var router, ok = this.children[name]
	if !ok {
		router, ok = this.regchildren[name]
	}
	return router, ok
}

// AddChild 添加子路由
func (this *BaseRouter) AddChild(router Router) bool {
	if router != nil {
		var name = strings.ToLower(router.Name())
		var children map[string]Router
		if router.Reg() {
			children = this.regchildren
		} else {
			children = this.children
		}
		var _, ok = children[name]
		if !ok {
			children[strings.ToLower(router.Name())] = router
			router.SetSuper(this)
			router.SetLevel(this.level + 1)
			return true
		}
	}
	return false
}

// AddChildren 批量添加添加子路由,如果已经存在同名路由,则添加失败
func (this *BaseRouter) AddChildren(routers ...Router) bool {
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
func (this *BaseRouter) AddBeforeFilter(filter RouterFilter) bool {
	if filter != nil {
		this.beforeFilters = append(this.beforeFilters, filter)
		return true
	}
	return false
}

// RemoveBeforeFilter 移除前置过滤器
func (this *BaseRouter) RemoveBeforeFilter(filter RouterFilter) bool {
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
func (this *BaseRouter) ExecBeforeFilter(context RouterContext) bool {
	for _, router := range this.beforeFilters {
		var goon = router.Filter(context)
		if !goon {
			return false
		}
	}
	return true
}

// AddAfterFilter 添加后置过滤器
func (this *BaseRouter) AddAfterFilter(filter RouterFilter) bool {
	if filter != nil {
		this.afterFilters = append(this.afterFilters, filter)
		return true
	}
	return false
}

// RemoveAfterFilter 移除后置过滤器
func (this *BaseRouter) RemoveAfterFilter(filter RouterFilter) bool {
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
func (this *BaseRouter) ExecAfterFilter(context RouterContext) bool {
	for _, router := range this.afterFilters {
		var goon = router.Filter(context)
		if !goon {
			return false
		}
	}
	return true
}
