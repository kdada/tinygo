package connector

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
	ErrorInvalidConnectorCreator Error = "ErrorInvalidConnectorCreator(N10000):无效的连接创建器"
	ErrorInvalidKind             Error = "ErrorInvalidKind(L10010):无效的连接类型(%s)"
	ErrorFailToStop              Error = "ErrorFailToStop(L10020):无法停止连接器(%s)"
	ErrorInvalidDispatcher       Error = "ErrorInvalidDispatcher(L10030):无效的Dispatcher,无法启动Connector(%s)"
	ErrorParamNotFound           Error = "ErrorParamNotFound(L10100):source中没有%s,无法创建https连接器"
)
