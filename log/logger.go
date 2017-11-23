package log

import (
	"container/list"
	"fmt"
	"path/filepath"
	"runtime"
	"strconv"
	"sync"
	"time"
)

// 日志记录器
type SimpleLogger struct {
	logLevel  LogLevel    //日志输出等级
	logList   *list.List  //日志链表
	logmu     *sync.Mutex //锁
	logWriter LogWriter   //日志写入
	closed    bool        //日志是否已经关闭
	async     bool        //日志是否异步输出
	skip      int         //日志输出文件信息时跳过的Caller数量
}

// NewSimpleLogger 创建日志记录器,默认同步模式
func NewSimpleLogger(logWriter LogWriter) *SimpleLogger {
	return &SimpleLogger{
		LogLevelDebug | LogLevelInfo | LogLevelWarn | LogLevelError | LogLevelFatal,
		list.New(),
		new(sync.Mutex),
		logWriter,
		false,
		false,
		2,
	}
}

// writeLog 写入日志
func (this *SimpleLogger) writeLog(info string, level LogLevel) {
	if !this.closed && level&this.logLevel > 0 {
		info = time.Now().Format("2006-01-02 15:04:05.000000 ") + info
		this.logmu.Lock()
		if this.async {
			this.logList.PushBack(info)
		} else {
			this.logWriter.Write(info)
		}
		this.logmu.Unlock()
	}
}

// file 返回调用Logger的函数位置
func (this *SimpleLogger) position(prefix string) string {
	if this.skip >= 2 {
		var _, file, line, ok = runtime.Caller(this.skip)
		if ok {
			return prefix + filepath.Base(file) + "(" + strconv.Itoa(line) + "):"
		}
	}
	return prefix
}

// Debug 写入调试信息
func (this *SimpleLogger) Debug(info ...interface{}) {
	this.writeLog(this.position("[Debug]")+fmt.Sprint(info), LogLevelDebug)
}

// Info 写入一般信息
func (this *SimpleLogger) Info(info ...interface{}) {
	this.writeLog(this.position("[Info]")+fmt.Sprint(info), LogLevelInfo)
}

// Warn 写入警告信息
func (this *SimpleLogger) Warn(info ...interface{}) {
	this.writeLog(this.position("[Warn]")+fmt.Sprint(info), LogLevelWarn)
}

// Error 写入错误信息
func (this *SimpleLogger) Error(info ...interface{}) {
	this.writeLog(this.position("[Error]")+fmt.Sprint(info), LogLevelError)
}

// Fatal 写入崩溃信息
func (this *SimpleLogger) Fatal(info ...interface{}) {
	this.writeLog(this.position("[Fatal]")+fmt.Sprint(info), LogLevelFatal)
}

// LogLevel 得到日志等级是否输出
func (this *SimpleLogger) LogLevelOutput(level LogLevel) bool {
	return this.logLevel&level > 0
}

// SetLogLevel 设置某个日志等级是否输出
func (this *SimpleLogger) SetLogLevelOutput(level LogLevel, output bool) {
	if output {
		this.logLevel |= level
	} else {
		this.logLevel &= ^level
	}
}

// Async 是否异步输出
func (this *SimpleLogger) Async() bool {
	return this.async
}

// SetAsync 设置是否异步输出
func (this *SimpleLogger) SetAsync(async bool) {
	var oldAsync = this.async
	this.async = async
	this.logWriter.SetAsync(this.async, this.logList, this.logmu)
	if oldAsync && !this.async && this.logList.Len() > 0 {
		//从异步切换回同步,将尚未异步输出的日志转换为同步输出
		var start *list.Element
		var length = 0
		this.logmu.Lock()
		if this.logList.Len() > 0 {
			start = this.logList.Front()
			length = this.logList.Len()
			this.logList.Init()
		}
		for i := 0; i < length; i++ {
			var v, ok = start.Value.(string)
			if ok {
				this.logWriter.Write(v)
			}
			start = start.Next()
		}
		this.logmu.Unlock()
	}
}

// LogPosition skip为跳过的Caller数量,skip小于2时关闭文件位置记录的功能
func (this *SimpleLogger) SetSkip(skip int) {
	this.skip = skip
}

// Closed 日志是否已关闭
func (this *SimpleLogger) Closed() bool {
	return this.closed
}

// Close 关闭日志 关闭后无法再使用
func (this *SimpleLogger) Close() {
	if !this.Closed() {
		this.closed = true
		this.logWriter.Close()
	}
}
