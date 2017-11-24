// Package log 包含日志相关工具
package log

import "sync"

// 日志接口 所有日志类型应该实现该接口
type Logger interface {
	// Debug 写入调试信息
	Debug(info ...interface{})
	// Info 写入一般信息
	Info(info ...interface{})
	// Warn 写入异常信息
	Warn(info ...interface{})
	// Error 写入错误信息
	Error(info ...interface{})
	// Fatal 写入崩溃信息
	Fatal(info ...interface{})
	// LogLevel 获取某个日志等级是否输出
	LogLevelOutput(level LogLevel) bool
	// SetLogLevelOutput 设置某个日志等级是否输出
	SetLogLevelOutput(level LogLevel, output bool)
	// Async 是否异步输出
	Async() bool
	// SetAsync 设置是否异步输出
	SetAsync(async bool)
	// Close 关闭日志 关闭后无法再进行写入操作
	Close()
	// Closed 日志是否关闭
	Closed() bool
}

// 日志创建器
//  source: 日志存储位置
type LoggerCreator func(source string) (Logger, error)

var (
	mu       sync.Mutex                       //互斥锁
	creators = make(map[string]LoggerCreator) //日志创建器映射
)

// NewLogger 创建一个新的Logger
//  kind:日志类型
func NewLogger(kind string, param string) (Logger, error) {
	var creator, ok = creators[kind]
	if !ok {
		return nil, ErrorInvalidKind.Format(kind).Error()
	}
	return creator(param)
}

// Register 注册LoggerCreator创建器
func Register(kind string, creator LoggerCreator) {
	if creator == nil {
		panic(ErrorInvalidLoggerCreator)
	}
	mu.Lock()
	defer mu.Unlock()
	creators[kind] = creator
}
