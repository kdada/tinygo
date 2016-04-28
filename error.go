package tinygo

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
	ErrorConfigNotCorrect     Error = "ErrorConfigNotCorrect(T10010):配置文件中%s的%s不正确"
	ErrorConnectorCreateFail  Error = "ErrorConnectorCreateFail(T10020):连接器(%s)创建失败,%s"
	ErrorRootRouterCreateFail Error = "ErrorRootRouterCreateFail(T10030):根路由(%s)创建失败,%s"
)
