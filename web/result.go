package web

import "net/http"

// 服务端状态码
type StatusCode int

const (
	//Http状态码
	StatusCodeOK               StatusCode = 200 //http正常返回结果
	StatusCodeMovedPermanently StatusCode = 301 //http永久转移
	StatusCodeMovedTemporarily StatusCode = 302 //http临时转移
	StatusCodeNotFound         StatusCode = 404 //http页面未找到
	//框架内部状态码
	StatusCodeRouterNotFound  StatusCode = iota + 10000 //路由未找到
	StatusCodeParamNotEnough                            //http参数不足
	StatusCodeParamNotCorrect                           //http参数不正确
	StatusCodePageNotFound                              //页面未找到
	//用户自定义状态码
	StatusCodeUserDefined StatusCode = 1000000 //用户自定义状态码起始码
)

// 可用于所有http方法的返回结果
type Result interface {
	// Code 状态码
	Code() StatusCode
	// Message 信息
	Message() string
	// Write 将Result的内容写入http.ResponseWriter
	Write(resp http.ResponseWriter)
}

// 可用于Get方法的返回结果
type GetResult Result

// 可用于Post方法的返回结果
type PostResult Result

// 可用于Put方法的返回结果
type PutResult Result

// 可用于Delete方法的返回结果
type DeleteResult Result

// 可用于Options方法的返回结果
type OptionsResult Result

// 可用于Head方法的返回结果
type HeadResult Result

// 可用于Trace方法的返回结果
type TraceResult Result

// 可用于Connect方法的返回结果
type ConnectResult Result

// 可用于Get和Post方法的返回结果
type GetPostResult Result
