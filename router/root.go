package router

// 根路由
// 根路由不占用路由层级,仅仅用于作为路由的根使用
// 可以认为根路由为 "/"
type RootRouter struct {
	SpaceRouter
}

// NewSpaceRouter 创建空间路由
// name:路由名称
func NewRootRouter() *SpaceRouter {
	var router = new(SpaceRouter)
	router.Init("Root")
	return router
}

// Check 检查当前路由器是否支持该请求
func (this *RootRouter) Check(routes []string) bool {
	return true
}

// SetSuper 设置父路由
func (this *RootRouter) SetSuper(super IRouter) {
	//根路由无法设置
}

// SetLevel 设置当前路由层级
func (this *RootRouter) SetLevel(level int) {
	//根路由无法设置
}
