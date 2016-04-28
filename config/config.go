// Package config 包含解析配置文件相关工具
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
	// Sections 返回全部节
	Sections() []Section
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
	parsersMu sync.Mutex                      //互斥锁
	parsers   = make(map[string]ConfigParser) //配置解析器
)

// NewConfig 创建一个新的Config
//  path:配置文件路径
func NewConfig(kind string, path string) (Config, error) {
	var data, err = ioutil.ReadFile(path)
	if err != nil {
		return nil, ErrorReadError.Format(path).Error()
	}
	return NewConfigWithContent(kind, data)
}

// NewConfigWithContent 创建一个新的Config
//  content:配置文件内容
func NewConfigWithContent(kind string, content []byte) (Config, error) {
	var parser, ok = parsers[kind]
	if !ok {
		return nil, ErrorInvalidConfigKind.Format(kind).Error()
	}
	return parser(content)
}

// Register 注册ConfigParser
func Register(kind string, parser ConfigParser) {
	if parser == nil {
		panic(ErrorInvalidConfigParser)
	}
	parsersMu.Lock()
	defer parsersMu.Unlock()
	parsers[kind] = parser
}
