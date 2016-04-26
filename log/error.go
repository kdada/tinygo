package log

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
	ErrorInvalidLoggerCreator  Error = "ErrorInvalidLoggerCreator(L10000):无效的日志创建器"
	ErrorInvalidKind           Error = "ErrorInvalidKind(L10010):无效的日志类型(%s)"
	ErrorInvalidLogWriter      Error = "ErrorInvalidLogWriter(L10011):无效的日志写入器"
	ErrorLogWriterInvalidParam Error = "ErrorLogWriterInvalidParam(L10020):SimpleLogWriter.SetAsync参数为nil"
	ErrorFailToCreatePath      Error = "ErrorFailToCreatePath(L10030):无法创建日志目录(%s)"
	ErrorInvalidParam          Error = "ErrorInvalidParam(L10040):无效的日志创建参数"
)
