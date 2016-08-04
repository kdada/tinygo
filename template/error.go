package template

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
	ErrorParamMustBeFunc    Error = "ErrorParamMustBeFunc(T10010):参数必须是函数"
	ErrorInvalidPartialView Error = "ErrorInvalidPartialView(T10020):无效的部分视图(%s),找不到指定名称(%s)的模板"
)
