package router

type Router interface {
	// Name 返回当前路由名称
	Name() string
	// Parent 返回当前父路由,每个Router只能有一个Parent
	Parent() Router
	// SetParent 设置当前路由父路由,当前路由必须是父路由的子路由
	SetParent(router Router) error
	// Named 返回当前是否使用Name进行路由匹配
	Named() bool
	// AddChild 添加子路由,Name不能重复,否则后者的会覆盖前者
	AddChild(router Router)
	// AddChildren 批量添加子路由,Name不能重复,否则后者的会覆盖前者
	AddChildren(routers []Router)
	// Child 返回指定名称的子路由
	Child(name string) (Router, bool)
	// RemoveChild 移除指定名称的路由,并返回该路由
	RemoveChild(name string) (Router, bool)
	// AddPreFilter 添加前置过滤器
	AddPreFilter(filter PreFilter) Router
	// RemovePreFilter 移除前置过滤器
	RemovePreFilter(filter PreFilter) bool
	// ExecPreFilter 执行前置过滤器
	ExecPreFilter(context RouterContext) bool
	// AddPostFilter 添加后置过滤器
	AddPostFilter(filter PostFilter) Router
	// RemovePostFilter 移除后置过滤器
	RemovePostFilter(filter PostFilter) bool
	// ExecPostFilter 执行后置过滤器
	ExecPostFilter(context RouterContext) bool
	// SetRouterExcutor 设置路由执行器生成方法
	SetRouterExcutorGenerator(RouterExcutorGenerator)
	// Match 匹配指定路由上下文,匹配成功则返回RouterExcutor
	Match(context RouterContext) (RouterExcutor, bool)
}

// 路由执行器生成器,每次应当返回一个全新的RouterExcutor实例
type RouterExcutorGenerator func() RouterExcutor

// 路由上下文
type RouterContext interface {
	// Segments 返回可匹配路由段
	Segments() []string
	// Match 匹配数量
	Match(count int)
	// Unmatch 失配数量
	Unmatch(count int)
	// Value 返回路由值
	Value(name string) string
	// SetValue 设置路由值
	SetValue(name string, value string)
	// Terminated 路由过程是否终止
	Terminated() bool
	// Terminate 终止路由过程,终止后该路由上下文将不会被继续路由并且不会被执行器执行
	Terminate()
	// Data 返回路由上下文携带的信息
	Data() interface{}
}

// 路由执行器
type RouterExcutor interface {
	// Router 返回生成RouterExcutor的路由
	Router() Router
	// SetRouter 设置生成RouterExcutor的路由
	SetRouter(router Router)
	// RouterContext 返回路由上下文
	RouterContext() RouterContext
	// SetRouterContext 设置路由上下文
	SetRouterContext(context RouterContext)
	// Excute 执行
	Excute() error
}

// 前置过滤器
type PreFilter interface {
	// Filter 过滤该请求
	// return:返回true表示继续处理,否则终止路由过程,后续的过滤器也不会执行
	Filter(context RouterContext) bool
}

// 后置过滤器
type PostFilter interface {
	// Filter 过滤该请求
	// return:返回true表示继续处理,否则终止路由过程,后续的过滤器也不会执行
	Filter(context RouterContext) bool
}
