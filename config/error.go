package config

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
	ErrorInvalidConfigParser   Error = "ErrorInvalidConfigParser(C10000):无效的配置解析器"
	ErrorInvalidKey            Error = "ErrorInvalidKey(C10010):无效的key(%s)"
	ErrorInvalidTypeConvertion Error = "ErrorInvalidTypeConvertion(C10020):无效的类型转换,string(%s)转换为%s"
	ErrorInvalidConfigKind     Error = "ErrorInvalidConfigKind(C10030):无效的配置类型(%s)"
	ErrorReadError             Error = "ErrorReadError(C10040):无法读取配置文件(%s)"
	ErrorNotMatch              Error = "ErrorNotMatch(C10050):未找到匹配的字符(%s)"
)
