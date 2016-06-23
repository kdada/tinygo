package web

// 服务端状态码
type StatusCode int

const (
	//Http状态码
	StatusCodeOK               StatusCode = 200 //http正常返回结果
	StatusCodeMovedPermanently StatusCode = 301 //http永久转移
	StatusCodeMovedTemporarily StatusCode = 302 //http临时转移
	StatusCodeNotFound         StatusCode = 404 //http页面未找到

	//框架内部状态码
	StatusCodeParamNotCorrect StatusCode = iota + 10000 //http参数不正确
	StatusCodePageNotFound                              //路由未找到

	//用户自定义状态码
	StatusCodeUserDefined StatusCode = 1000000 //用户自定义状态码起始码
)
