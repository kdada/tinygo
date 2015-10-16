package tinygo

import (
	"fmt"
	"path/filepath"
	"runtime"
)

// Debug 输出调试信息,包含调用位置
func Debug(logs ...interface{}) {
	var _, file, line, _ = runtime.Caller(1)
	file, _ = filepath.Rel(tinyConfig.path, file)
	fmt.Print(file, "[", line, "] ")
	fmt.Println(logs...)
}
