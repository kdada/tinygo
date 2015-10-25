// Package config 实现了一个ini配置文件解析器
package config

import (
	"io/ioutil"
	"sync"
)

// 配置信息存储接口
type Config interface {
	// GlobalSection 获取全局配置段
	GlobalSection() Section
	// Section 根据name获取指定名称的配置段
	Section(name string) (Section, bool)
}

// 配置信息段
type Section interface {
	// Name 配置段名称
	Name() string
	// String 获取字符串
	String(key string) (string, error)
	// Int 获取整数
	Int(key string) (int64, error)
	// Bool 获取布尔值
	Bool(key string) (bool, error)
	// Float 获取浮点值
	Float(key string) (float64, error)
}

// 配置解析方法
type ConfigParser func([]byte) (Config, error)

var (
	parsersMu sync.Mutex                          //生成器互斥锁
	parsers   = make(map[ConfigType]ConfigParser) //配置解析器
)

// NewConfig 创建一个新的Config
//  path:配置文件路径
func NewConfig(kind ConfigType, path string) (Config, error) {
	var data, err = ioutil.ReadFile(path)
	if err != nil {
		return nil, ConfigErrorReadError.Format(path).Error()
	}
	var parser, ok = parsers[kind]
	if !ok {
		return nil, ConfigErrorInvalidConfigKind.Format(kind).Error()
	}
	return parser(data)
}

// registerProviderCreator 注册InjectorProvider创建器
func registerConfigParser(kind ConfigType, parser ConfigParser) error {
	parsersMu.Lock()
	defer parsersMu.Unlock()
	parsers[kind] = parser
	return nil
}
