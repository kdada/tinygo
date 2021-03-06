package web

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
	ErrorNotStructPtr            Error = "ErrorNotStructPtr(W10030):(%s)不是结构体指针类型"
	ErrorNoSpecificMethod        Error = "ErrorNoSpecificMethod(W10031):控制器(%s)不存在指定的方法(%s)"
	ErrorFirstReturnMustBeResult Error = "ErrorFirstReturnMustBeResult(W10040):第一个返回值类型(%s)不符合web.Result接口"
	ErrorNoReturn                Error = "ErrorNoReturn(W10041):函数(%s)至少拥有一个返回值并且第一个返回值必须符合web.Result类型"

	ErrorParamNotExist  Error = "ErrorParamNotExist(W10100):参数(%s)不存在"
	ErrorRouterNotFound Error = "ErrorRouterNotFound(W10110):路由(%s)不存在"
	ErrorInvalidContext Error = "ErrorInvalidContext(W10120):无效的上下文(%s),无法转换为web.Context"

	ErrorInvalidWriter      Error = "ErrorInvalidWriter(W10200):无效的http写入器"
	ErrorInvalidPartialView Error = "ErrorInvalidPartialView(W10300):无效的部分视图(%s),找不到指定名称(%s)的模板"
	ErrorInvalidKey         Error = "ErrorInvalidKey(W10400):无效的Key(%s)"
	ErrorParamMustBeFunc    Error = "ErrorParamMustBeFunc(W10500):参数必须是函数"
)
