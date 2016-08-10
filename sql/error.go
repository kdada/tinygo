package sql

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
	ErrorInvalidParamType Error = "Q10000:ErrorInvalidParamType,无效的输入类型(%s)"
	ErrorBeginFailed      Error = "Q10010:ErrorBeginFailed,当前DB已经打开了一个事务"
)
