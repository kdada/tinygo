package log

//日志级别
type LogLevel uint

const (
	LogLevelDebug LogLevel = 1 << iota //调试信息等级
	LogLevelInfo                       //输出信息等级
	LogLevelWarn                       //警告信息等级
	LogLevelError                      //错误信息等级
	LogLevelFatal                      //崩溃信息等级
)
