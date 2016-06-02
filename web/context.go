package web

import (
	"github.com/kdada/tinygo/connector"
	"github.com/kdada/tinygo/router"
)

type Context struct {
	router.BaseContext
	Data        map[string]string      //路由信息
	HttpContext *connector.HttpContext //http上下文
}

// NewContext 创建上下文信息
func NewContext(segments []string, context *connector.HttpContext) *Context {
	var method = context.Request.Method
	var c = new(Context)
	c.Data = make(map[string]string, 1)
	c.Segs = append(segments, method)
	c.HttpContext = context
	return c
}

// Value 返回路由值
func (this *Context) Value(name string) (string, bool) {
	var value, ok = this.Data[name]
	return value, ok
}

// SetValue 设置路由值
func (this *Context) SetValue(name string, value string) {
	this.Data[name] = value
}
