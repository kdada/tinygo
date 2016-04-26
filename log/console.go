package log

import (
	"fmt"
)

// 控制台写入器
type ConsoleWriter struct {
}

// NewConsoleWriter 创建控制台写入器
func NewConsoleWriter() *ConsoleWriter {
	return new(ConsoleWriter)
}

// Write 日志写入
func (this *ConsoleWriter) Write(log string) {
	fmt.Println(log)
}

// Close 关闭写入器
func (this *ConsoleWriter) Close() {

}

// ConsoleLoggerCreator
func ConsoleLoggerCreator(param interface{}) (Logger, error) {
	return NewSimpleLogger(NewSimpleLogWriter(NewConsoleWriter())), nil
}
