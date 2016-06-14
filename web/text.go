package web

import (
	"fmt"
	"runtime"
)

// Info 生成信息,包含调用当前方法的文件位置
func Info(logs ...interface{}) string {
	var _, file, line, _ = runtime.Caller(1)
	return fmt.Sprint(file, ":", line, "\n\t") + fmt.Sprintln(logs...)
}

// Output 输出信息,包含调用当前方法的文件位置
func Output(logs ...interface{}) {
	var _, file, line, _ = runtime.Caller(1)
	fmt.Print(file, ":", line, "\n\t")
	fmt.Println(logs...)
}
