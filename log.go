package tinygo

import (
	"fmt"
	"log"
	"path/filepath"
	"runtime"
)

// Debug 输出调试信息,包含调用位置
// 如果当前处于发布模式,则不输出logs
func Debug(logs ...interface{}) {
	if !IsRelease() {
		var info = "[DEBUG] " + outputLineInfo()
		var allInfo = info + fmt.Sprintln(logs...)
		log.Print(allInfo)
	}
}

// Log 输出信息,包含调用位置
// 即使当前处于发布模式,也会输出logs
func Log(logs ...interface{}) {
	var info = "[LOG] " + outputLineInfo()
	var allInfo = info + fmt.Sprintln(logs...)
	log.Print(allInfo)
}

// Error 输出错误信息,包含调用位置
// 即使当前处于发布模式,也会输出logs
func Error(logs ...interface{}) {
	var info = "[ERROR] " + outputLineInfo()
	var allInfo = info + fmt.Sprintln(logs...)
	log.Print(allInfo)
}

// outputLineInfo 生成行信息
func outputLineInfo() string {
	var _, file, line, _ = runtime.Caller(2)
	var _, fileName = filepath.Split(file)
	return fmt.Sprint(fileName, ":", line, " ")
}
