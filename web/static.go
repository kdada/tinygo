package web

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/kdada/tinygo/router"
)

// 静态文件执行器
type StaticExecutor struct {
	router.BaseRouterExecutor
	path string
}

// NewStaticExecutor 创建静态文件执行器
func NewStaticExecutor(path string) *StaticExecutor {
	var se = new(StaticExecutor)
	se.path = path
	return se
}

// Excute 执行
func (this *StaticExecutor) Execute() (interface{}, error) {
	var context, ok = this.Context.(*Context)
	if ok {
		context.End = this.End
		return this.FilterExecute(func() (interface{}, error) {
			var result Result = nil
			if context.HttpContext.Request.Method == "GET" {
				//返回文件
				var pathSegs = context.Segments()
				var containDotDot = false
				for _, s := range pathSegs {
					if strings.Contains(s, "..") {
						containDotDot = true
						break
					}
				}
				if !containDotDot {
					var r = this.Router()
					var count = 0
					for r != nil {
						r = r.Parent()
						count++
					}
					var filePath = filepath.Join(this.path, strings.Join(pathSegs[count:len(pathSegs)-1], "/"))
					if !context.Processor.Config.List {
						var f, err = os.Stat(filePath)
						if err != nil || f.IsDir() {
							result = context.NotFound()
						} else {
							result = context.File(filePath)
						}
					} else {
						result = context.File(filePath)
					}

				}
			}
			if result == nil {
				result = context.NotFound()
			}
			return result, nil
		})
	}
	return nil, ErrorInvalidContext.Format(reflect.TypeOf(this.Context).String()).Error()
}

// 文件执行器,用于返回特定文件
type FileExecutor struct {
	router.BaseRouterExecutor
	path string
}

// NewFileExecutor 创建文件执行器
func NewFileExecutor(path string) *FileExecutor {
	var fe = new(FileExecutor)
	fe.path = path
	return fe
}

// Excute 执行
func (this *FileExecutor) Execute() (interface{}, error) {
	var context, ok = this.Context.(*Context)
	if ok {
		context.End = this.End
		return this.FilterExecute(func() (interface{}, error) {
			return context.File(this.path), nil
		})
	}
	return nil, ErrorInvalidContext.Format(reflect.TypeOf(this.Context).String()).Error()
}
