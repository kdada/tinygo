package info

import (
	"errors"
	"fmt"
)

// 配置错误信息
type TinyGoError string

// 错误码
const (
	TinyGoErrorConfigError TinyGoError = "T10010:ConfigErrorInvalidKey,无效的key(%s)"
)

// Format 格式化错误信息并生成新的错误信息
func (this TinyGoError) Format(data ...interface{}) TinyGoError {
	return TinyGoError(fmt.Sprintf(string(this), data...))
}

// Error 生成error类型
func (this TinyGoError) Error() error {
	return errors.New(string(this))
}
