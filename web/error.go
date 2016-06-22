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
	ErrorNotMethod               Error = "ErrorNotMethod(W10010):%s不是函数"
	ErrorParamNotPtr             Error = "ErrorParamNotPtr(W10020):函数(%s)的参数类型(%s)不是结构体指针类型"
	ErrorNotStructPtr            Error = "ErrorNotStructPtr(W10030):(%s)不是结构体指针类型"
	ErrorFirstReturnMustBeResult Error = "ErrorFirstReturnMustBeResult(W10040):第一个返回值类型(%s)不符合web.Result接口"
	ErrorNoReturn                Error = "ErrorNoReturn(W10041):函数(%s)至少拥有一个返回值并且第一个返回值必须符合web.Result类型"
	ErrorParamNotExist           Error = "ErrorParamNotExist(W10100):参数(%s)不存在"
	ErrorInvalidWriter           Error = "ErrorInvalidWriter(W10200):无效的http写入器"
)
