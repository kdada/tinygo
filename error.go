package tinygo

import (
	"errors"
	"fmt"
)

// 配置错误信息
type TinyGoError string

// 错误码
const (
	TinyGoErrorParamNotFoundError TinyGoError = "T10110:TinyGoErrorParamNotFoundError,不存在的key(%s)"
)

// Format 格式化错误信息并生成新的错误信息
func (this TinyGoError) Format(data ...interface{}) TinyGoError {
	return TinyGoError(fmt.Sprintf(string(this), data...))
}

// Error 生成error类型
func (this TinyGoError) Error() error {
	return errors.New(string(this))
}
