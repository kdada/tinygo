package log

import (
	"container/list"
	"sync"
	"sync/atomic"
	"time"
)

// 日志写入类型接口
type LogWriter interface {
	// Write [同步]日志写入
	Write(log string)
	// SetAsync 设置写入模式
	SetAsync(async bool, logList *list.List, mu *sync.Mutex)
	// Close 关闭写入器
	Close()
}

// 写入类型接口
type Writer interface {
	// Write 日志写入
	Write(log string)
	// Close 关闭写入器
	Close()
}

// 日志写入器
type SimpleLogWriter struct {
	logList *list.List  // 日志列表
	logmu   *sync.Mutex // 日志列表锁
	async   bool        // 当前是否处于异步模式
	counter int32       // 协程计数器
	closed  bool        // 是否已经停止
	writer  Writer      // 写入器
}

// NewSimpleLogWriter 创建日志写入器
func NewSimpleLogWriter(writer Writer) *SimpleLogWriter {
	return &SimpleLogWriter{
		nil,
		nil,
		false,
		0,
		false,
		writer,
	}
}

// Write [同步]日志写入
func (this *SimpleLogWriter) Write(log string) {
	if !this.closed && !this.async && this.writer != nil {
		this.writer.Write(log)
	}
}

// SetAsync 设置是否异步输出,仅当之前为非异步并且设置为异步时,logList和mu有效
func (this *SimpleLogWriter) SetAsync(async bool, logList *list.List, mu *sync.Mutex) {
	if !this.async && async {
		if logList == nil || mu == nil {
			panic(ErrorLogWriterInvalidParam)
		}
		this.logList = logList
		this.logmu = mu
		go func() {
			var counter = atomic.AddInt32(&this.counter, 1)
			for this.async && counter == this.counter && !this.closed {
				if this.logList.Len() > 0 {
					var start *list.Element
					var length = 0
					this.logmu.Lock()
					if this.logList.Len() > 0 {
						start = this.logList.Front()
						length = this.logList.Len()
						this.logList.Init()
					}
					this.logmu.Unlock()
					for i := 0; i < length; i++ {
						var v, ok = start.Value.(string)
						if ok {
							this.writer.Write(v)
						}
						start = start.Next()
					}
				} else {
					time.Sleep(50 * time.Millisecond)
				}
			}
		}()
	}
	this.async = async
}

// Close 关闭日志写入器
func (this *SimpleLogWriter) Close() {
	this.async = false
	this.closed = true
	this.writer.Close()
}
