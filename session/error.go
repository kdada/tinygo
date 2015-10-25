package session

import (
	"errors"
	"fmt"
)

// 错误信息
type SessionError string

// 错误码
const (
	SessionErrorInvalidSessionKind SessionError = "S10010:SessionErrorInvalidSessionKind,无效的Session Kind(%s)"
)

// Format 格式化错误信息并生成新的错误信息
func (this SessionError) Format(data ...interface{}) SessionError {
	return SessionError(fmt.Sprintf(string(this), data...))
}

// Error 生成error类型
func (this SessionError) Error() error {
	return errors.New(string(this))
}
