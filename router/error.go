package router

import (
	"errors"
	"fmt"
)

// 错误信息
type Error string

// Format 格式化错误信息并生成新的错误信息
func (this Error) Format(data ...interface{}) Error {
	return Error(fmt.Sprintf(string(this), data...))
}

// Error 生成error类型
func (this Error) Error() error {
	return errors.New(string(this))
}

// String 返回错误字符串描述
func (this Error) String() string {
	return string(this)
}

// 错误码
const (
	ErrorInvalidParentRouter  Error = "ErrorInvalidParentRouter(R10010):指定路由不是当前路由的父路由,设置父路由失败"
	ErrorNamedRouterNoChecker Error = "ErrorNamedRouterNoChecker(R10020):名称路由不能设置检查器"
	ErrorRegexpParseError     Error = "ErrorRegexpParseError(R10030),无效的正则表达式(%s)"
	ErrorRegexpNoneError      Error = "ErrorRegexpNoneError(R10031),字符串(%s)不包含正则表达式"
	ErrorRegexpFormatError    Error = "ErrorRegexpFormatError(R10032),字符串(%s)格式错误"
	ErrorRegexpNotMatchError  Error = "ErrorRegexpNotMatchError(R10033),正则表达式匹配失败"
	ErrorExecutorDoNothing    Error = "ErrorExecutorDoNothing(R10050),空执行器错误,该执行器没有执行任何内容"
	ErrorInvalidKind          Error = "ErrorInvalidKind(R10060),无效的路由类型(%s)"
	ErrorInvalidRouterCreator Error = "ErrorInvalidRouterCreator(R10070),无效的路由创建器"
	ErrorInvalidMatchParam    Error = "ErrorInvalidMatchParam(R10080),无效的match参数(%s),期望参数类型为%s"
)
