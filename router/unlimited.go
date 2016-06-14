package router

// 无限路由
type UnlimitedRouter struct {
	self              Router                 //路由自身(所有继承当前路由的路由在创建时需要设置self=this)
	parent            Router                 //父路由
	name              string                 //当前路由名称
	preFilters        []PreFilter            //在子路由处理之前执行的过滤器
	postFilters       []PostFilter           //在子路由处理之后执行的过滤器
	executorGenerator RouterExcutorGenerator //路由执行器生成器
}

// NewUnlimitedRouter 创建无限路由
//  name:无限路由名称
//  match:无限路由不需要匹配规则,该参数无效
func NewUnlimitedRouter(name string, match interface{}) (Router, error) {
	var r = new(UnlimitedRouter)
	r.name = name
	r.preFilters = make([]PreFilter, 0)
	r.postFilters = make([]PostFilter, 0)
	r.self = r
	return r, nil
}

// Name 返回当前路由名称
func (this *UnlimitedRouter) Name() string {
	return this.name
}

// MatchString 返回当前路由用于进行匹配的字符串
func (this *UnlimitedRouter) MatchString() string {
	return ""
}

// Parent 返回当前父路由,每个Router只能有一个Parent
func (this *UnlimitedRouter) Parent() Router {
	return this.parent
}

// SetParent 设置当前路由父路由,当前路由必须是父路由的子路由
func (this *UnlimitedRouter) SetParent(router Router) error {
	var r, ok = router.Child(this.name)
	var x, ok2 = r.(*UnlimitedRouter)
	if ok && ok2 && x == this {
		if this.parent != nil && this.parent != router {
			this.parent.RemoveChild(this.name)
		}
		this.parent = router
		return nil
	}
	return ErrorInvalidParentRouter.Error()
}

// Normal 返回当前路由是否为通常路由,通常路由可以使用MatchString()返回的字符串进行直接匹配
func (this *UnlimitedRouter) Normal() bool {
	return false
}

// AddChild 无限路由不能添加子路由
func (this *UnlimitedRouter) AddChild(router Router) {

}

// AddChildren 无限路由不能添加子路由
func (this *UnlimitedRouter) AddChildren(routers []Router) {

}

// Child 无限路由没有子路由
func (this *UnlimitedRouter) Child(name string) (Router, bool) {
	return nil, false
}

// Children 无限路由没有子路由
func (this *UnlimitedRouter) Children() []Router {
	return []Router{}
}

// RemoveChild 无限路由没有子路由
func (this *UnlimitedRouter) RemoveChild(name string) (Router, bool) {
	return nil, false
}

// AddPreFilter 添加前置过滤器
func (this *UnlimitedRouter) AddPreFilter(filter PreFilter) Router {
	if filter != nil {
		this.preFilters = append(this.preFilters, filter)
	}
	return this.self
}

// RemovePreFilter 移除前置过滤器
func (this *UnlimitedRouter) RemovePreFilter(filter PreFilter) bool {
	for index, child := range this.preFilters {
		if child == filter {
			this.preFilters = append(this.preFilters[:index], this.preFilters[index+1:]...)
			return true
		}
	}
	return false
}

// ExecPreFilter 执行前置过滤器
func (this *UnlimitedRouter) ExecPreFilter(context RouterContext) bool {
	for _, router := range this.preFilters {
		var goon = router.Filter(context)
		if !goon {
			return false
		}
	}
	return true
}

// AddPostFilter 添加后置过滤器
func (this *UnlimitedRouter) AddPostFilter(filter PostFilter) Router {
	if filter != nil {
		this.postFilters = append(this.postFilters, filter)
	}
	return this.self
}

// RemovePostFilter 移除后置过滤器
func (this *UnlimitedRouter) RemovePostFilter(filter PostFilter) bool {
	for index, child := range this.postFilters {
		if child == filter {
			this.postFilters = append(this.postFilters[:index], this.postFilters[index+1:]...)
			return true
		}
	}
	return false
}

// ExecPostFilter 执行后置过滤器
func (this *UnlimitedRouter) ExecPostFilter(context RouterContext, result interface{}) bool {
	for _, router := range this.postFilters {
		var goon = router.Filter(context, result)
		if !goon {
			return false
		}
	}
	return true
}

// SetRouterExcutor 设置路由执行器生成方法
func (this *UnlimitedRouter) SetRouterExcutorGenerator(reg RouterExcutorGenerator) {
	this.executorGenerator = reg
}

// Match 匹配指定路由上下文,匹配成功则返回RouterExcutor
func (this *UnlimitedRouter) Match(context RouterContext) (RouterExcutor, bool) {
	if this.executorGenerator != nil {
		var executor = this.executorGenerator()
		executor.SetRouter(this)
		executor.SetRouterContext(context)
		return executor, true
	}
	//匹配失败
	return nil, false
}
