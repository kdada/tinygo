package router

type internalController interface {
	// SetContext 设置请求上下文环境
	SetContext(context RouterContext)
	// SetRouter 设置使用当前控制器的路由
	SetRouter(router Router)
}

// 控制器路由信息接口
type ControllerRouter interface {
	// 返回控制器路由信息
	//  return:方法名,路由名称,http方法,后缀段
	Info() (string, string, string, []string)
}

// 控制器接口
type Controller interface {
	internalController
	// Routers 返回当前控制器可以使用的方法路由信息
	Routers() []interface{}
}

// Restful控制器接口
type RestfulController interface {
	internalController
	// Get HTTP GET对应方法
	Get()
	// Post HTTP POST对应方法
	Post()
	// Put HTTP PUT对应方法
	Put()
	// Delete HTTP DELETE对应方法
	Delete()
}
