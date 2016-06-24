package web

// 服务端状态码
type StatusCode int

const (
	//Http状态码
	StatusCodeOK               StatusCode = 200 //http正常返回结果
	StatusCodeMovedPermanently StatusCode = 301 //http永久转移
	StatusCodeMovedTemporarily StatusCode = 302 //http临时转移
	StatusCodeNotFound         StatusCode = 404 //http页面未找到

	//框架内部状态码(功能)
	StatusCodeRedispatch StatusCode = iota + 10000 //路由重新分发状态,接收该状态后需要将当前请求重新分发
	//框架内部状态码(错误)
	StatusCodeParamNotCorrect //http参数不正确
	StatusCodePageNotFound    //路由未找到

	//用户自定义状态码
	StatusCodeUserDefined StatusCode = 1000000 //用户自定义状态码起始码
)
