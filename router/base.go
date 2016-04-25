package router

import "regexp"

// 基础路由
type BaseRouter struct {
	UnlimitedRouter
	reg    bool
	regexp *regexp.Regexp //正则表达式
	keys   []string       //可提取keys
}

// Named 返回当前是否使用Name进行路由匹配
func (this *BaseRouter) Named() bool {
	return !this.reg
}

// AddChild 添加子路由
func (this *BaseRouter) AddChild(router Router) {
	var _, ok = this.Child(router.Name())
	if ok {
		this.RemoveChild(router.Name())
	}
	if router.Named() {
		this.namedChildren[router.Name()] = router
	} else {
		this.unnamedChildren[router.Name()] = router
	}
}

// AddChildren 批量添加子路由
func (this *BaseRouter) AddChildren(routers []Router) {
	for _, v := range routers {
		this.AddChild(v)
	}
}

// Child 返回指定名称的子路由
func (this *BaseRouter) Child(name string) (Router, bool) {
	var r, ok = this.namedChildren[name]
	if ok {
		return r, ok
	}
	r, ok = this.unnamedChildren[name]
	return r, ok
}

// RemoveChild 移除指定名称的路由,并返回该路由
func (this *BaseRouter) RemoveChild(name string) (Router, bool) {

	var r, ok = this.namedChildren[name]
	if ok {
		delete(this.namedChildren, name)
		return r, ok
	}
	r, ok = this.unnamedChildren[name]
	if ok {
		delete(this.unnamedChildren, name)
		return r, ok
	}
	return nil, false
}

// Match 匹配指定路由上下文,匹配成功则返回RouterExcutor
func (this *BaseRouter) Match(context RouterContext) (RouterExcutor, bool) {
	var segs = context.Segments()
	var data [][]string
	if !this.Named() {
		data = this.regexp.FindAllStringSubmatch(segs[0], 1)
		if len(data) <= 0 {
			return nil, false
		}
	}
	context.Match(1)
	defer context.Unmatch(1)

	segs = context.Segments()
	var executor RouterExcutor
	if len(segs) <= 0 && this.executorGenerator == nil {
		//无效匹配
		context.Terminate()
		return nil, false
	}

	if len(segs) > 0 {
		//路由传递
		var name = segs[0]
		var r, ok = this.namedChildren[name]
		if ok {
			executor, ok = r.Match(context)
		}
		if !context.Terminated() && !ok {
			for _, v := range this.unnamedChildren {
				executor, ok = v.Match(context)
				if ok {
					break
				}
				if context.Terminated() {
					return nil, false
				}
			}
		}
		if !ok {
			return nil, false
		}
	}
	//匹配
	if executor == nil && this.executorGenerator != nil {
		executor = this.executorGenerator()
		executor.SetRouter(this)
		executor.SetRouterContext(context)
	}
	if len(data) >= 1 {
		var values = data[0][1:]
		for i, v := range values {
			context.SetValue(this.keys[i], v)
		}
	}
	return executor, true
}
