package tinygo

import (
	"fmt"
	"path/filepath"
	"runtime"
)

// Debug 输出调试信息,包含调用位置
// 如果当前处于发布模式,则不输出任何信息
func Debug(logs ...interface{}) {
	if !IsRelease() {
		var _, file, line, _ = runtime.Caller(1)
		file, _ = filepath.Rel(tinyConfig.path, file)
		fmt.Print(file, "[", line, "] ")
		fmt.Println(logs...)
	}
}

// Log 输出信息,包含调用位置
func Log(logs ...interface{}) {
	var _, file, line, _ = runtime.Caller(1)
	file, _ = filepath.Rel(tinyConfig.path, file)
	fmt.Print(file, "[", line, "] ")
	fmt.Println(logs...)
}
