package router

import "sync"

type Router interface {
	// Name 返回当前路由名称
	Name() string
	// MatchString 返回当前路由用于进行匹配的字符串
	MatchString() string
	// Parent 返回当前父路由,每个Router只能有一个Parent
	Parent() Router
	// SetParent 设置当前路由父路由,当前路由必须是父路由的子路由
	SetParent(router Router) error
	// Normal 返回当前路由是否为通常路由,通常路由可以使用MatchString()返回的字符串进行直接匹配
	Normal() bool
	// AddChild 添加子路由,Name相同的路由自动合并
	AddChild(router Router)
	// AddChildren 批量添加子路由,Name相同的路由自动合并
	AddChildren(routers []Router)
	// Child 返回指定名称的子路由
	Child(name string) (Router, bool)
	// Children 返回所有子路由
	Children() []Router
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
	ExecPostFilter(context RouterContext, result interface{}) bool
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
	Value(name string) (string, bool)
	// SetValue 设置路由值
	SetValue(name string, value string)
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
	// Execute 执行,并返回相应结果
	Execute() (interface{}, error)
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
	Filter(context RouterContext, result interface{}) bool
}

// 路由创建器
//  name: 路由名称
//  match: 用于进行匹配的内容,必须是指定路由所需要的内容
type RouterCreator func(name string, match interface{}) (Router, error)

var (
	mu       sync.Mutex                       //互斥锁
	creators = make(map[string]RouterCreator) //创建器映射
)

// NewRouter 创建一个新的Router
//  kind:路由类型
func NewRouter(kind string, name string, match interface{}) (Router, error) {
	var creator, ok = creators[kind]
	if !ok {
		return nil, ErrorInvalidKind.Format(kind).Error()
	}
	return creator(name, match)
}

// Register 注册RouterCreator创建器
func Register(kind string, creator RouterCreator) {
	if creator == nil {
		panic(ErrorInvalidRouterCreator)
	}
	mu.Lock()
	defer mu.Unlock()
	creators[kind] = creator
}
