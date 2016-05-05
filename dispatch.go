package tinygo

import "github.com/kdada/tinygo/router"

type Context struct {
	router.BaseContext
	routerInfo map[string]string
	data       interface{}
}

// NewContext 创建上下文信息
func NewContext(segments []string, data interface{}) *Context {
	var context = new(Context)
	context.routerInfo = make(map[string]string, 1)
	context.Segs = segments
	context.data = data
	return context
}

// Value 返回路由值
func (this *Context) Value(name string) (string, bool) {
	var value, ok = this.routerInfo[name]
	return value, ok
}

// SetValue 设置路由值
func (this *Context) SetValue(name string, value string) {
	this.routerInfo[name] = value
}

// Data 返回路由上下文携带的信息
func (this *Context) Data() interface{} {
	return this.data
}

// Dispatcher 调度器,用于协调连接器和路由
type Dispatcher struct {
	Root router.Router
}

// Dispatch 分发
//  segments:用于进行分发的路径段信息
//  data:连接携带的数据
func (this *Dispatcher) Dispatch(segments []string, data interface{}) {
	var context = NewContext(segments, data)
	var executor, ok = this.Root.Match(context)
	if ok {
		var err = executor.Execute()
		if err != nil {

		}
	} else {
	}
}
