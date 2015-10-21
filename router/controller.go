package router

//控制器接口
type Controller interface {
	// SetContext 设置请求上下文环境
	SetContext(context RouterContext)
	// SetRouter 设置使用当前控制器的路由
	SetRouter(router Router)
	// File 返回文件
	File(path string)
	// Json 返回json格式的数据
	Json(value interface{})
	// Xml 返回Xml格式的数据
	Xml(value interface{})
	// Api 根据设置返回Json或Xml
	Api(value interface{})
	// View 返回视图页面
	View(path string, data ...interface{})
	// SimpleView 返回 控制器名(不含Controller)/方法名.html 页面
	SimpleView(data ...interface{})
	// PartialView 返回 控制器名(不含Controller)/方法名.html 页面无视layout设置
	PartialView(path string, data ...interface{})
	// PartialView 返回 控制器名(不含Controller)/方法名.html 页面无视layout设置
	SimplePartialView(data ...interface{})
	// HttpNotFound 返回404
	HttpNotFound()
	// ParseParams 将参数解析到结构体中
	// params:结构体指针数组,参数必须是结构体指针
	ParseParams(params ...interface{})
	// RedirectMethod [302] 重定向到当前控制器的方法
	// method:方法名
	// params:要传递的参数(这些参数将作为query string传递)
	RedirectMethod(method string, params ...interface{})
	// Redirect [302] 重定向到指定控制器的指定方法
	// controller:控制器名,该控制器必须与当前控制器处于同一个SpaceRouter中
	// method:方法名
	// params:要传递的参数(这些参数将作为query string传递)
	Redirect(controller string, method string, params ...interface{})
	// Redirect [302] 重定向到指定url
	RedirectUrl(url string)
	// RedirectPermanently [301] 永久重定向到指定控制器的方法
	// controller:控制器名,该控制器必须与当前控制器处于同一个SpaceRouter中
	// method:方法名
	// params:要传递的参数(这些参数将作为query string传递)
	RedirectPermanently(controller string, method string, params ...interface{})
	// RedirectUrlPermanently [301] 永久重定向到指定url
	RedirectUrlPermanently(url string)
	// Routers 返回当前控制器可以使用的方法路由信息
	Routers() []RouterInfo
}

//Restful控制器接口
type RestfulController interface {
	Controller
	// Get HTTP GET对应方法
	Get()
	// Post HTTP POST对应方法
	Post()
	// Put HTTP PUT对应方法
	Put()
	// Delete HTTP DELETE对应方法
	Delete()
}
