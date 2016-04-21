package log

import (
	"fmt"
	"os"
	"testing"
	"time"
)

var consoleLogger Logger

func TestMain(m *testing.M) {
	var err error
	consoleLogger, err = NewLogger("console")
	if err != nil {
		fmt.Println(err)
		return
	}
	fileLogger, err = NewLogger("file")
	fileLogger.SetAsync(true)
	if err != nil {
		fmt.Println(err)
		return
	}
	os.Exit(m.Run())
}

//测试输出
func TestConsoleLoggerOutput(t *testing.T) {
	var logger = consoleLogger
	var log = "T:测试测试"
	logger.Debug(log)
	logger.Info(log)
	logger.Warn(log)
	logger.Error(log)
	logger.Fatal(log)

	time.Sleep(1 * time.Second)
}

//性能测试
func BenchmarkConsoleLoggerOutput(b *testing.B) {
	var logger = consoleLogger
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			var log = "B:测试测试"
			logger.Debug(log)
		}
	})
}
