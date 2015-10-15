package config

import (
	"errors"
	"fmt"
)

// 配置错误信息
type ConfigError string

// 错误码
const (
	ConfigErrorInvalidKey            ConfigError = "C10010:ConfigErrorInvalidKey,无效的key(%s)"
	ConfigErrorInvalidTypeConvertion ConfigError = "C10020:ConfigErrorInvalidTypeConvertion,无效的类型转换,string(%s)转换为%s"
	ConfigErrorInvalidConfigKind     ConfigError = "C10030:ConfigErrorInvalidConfigKind,无效的配置类型(%s)"
	ConfigErrorReadError             ConfigError = "C10040:ConfigErrorReadError,无法读取配置文件(%s)"
	ConfigErrorNotMatch              ConfigError = "C10050:ConfigErrorNotMatch,未找到匹配的字符(%s)"
)

// Format 格式化错误信息并生成新的错误信息
func (this ConfigError) Format(data ...interface{}) ConfigError {
	return ConfigError(fmt.Sprintf(string(this), data...))
}

// Error 生成error类型
func (this ConfigError) Error() error {
	return errors.New(string(this))
}
