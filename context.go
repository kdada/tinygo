package tinygo

import (
	"net/http"

	"github.com/kdada/tinygo/router"
	"github.com/kdada/tinygo/session"
)

// 路由环境
type HttpContext struct {
	urlParts       []string               //Url分段信息,每段一个字符串
	request        *http.Request          //http请求
	responseWriter http.ResponseWriter    //http响应
	session        session.Session        //http会话
	Params         map[string]string      //http参数,包含url,query,form的所有参数
	parsed         bool                   //存储参数是否已经解析过
	routers        []router.Router        //分派成功的路由链
	executor       router.ContextExecutor //存储最终执行Context的执行器
}

// Method 返回Http方法
func (this *HttpContext) Method() string {
	return this.request.Method
}

// ResponseWriter 返回ResponseWriter
func (this *HttpContext) ResponseWriter() http.ResponseWriter {
	return this.responseWriter
}

// Request 返回Request
func (this *HttpContext) Request() *http.Request {
	return this.request
}

// RouterParts 返回路由段
func (this *HttpContext) RouterParts() []string {
	return this.urlParts
}

// SetRouterParts 设置路由段
func (this *HttpContext) SetRouterParts(parts []string) {
	this.urlParts = parts
}

// AddParams 添加路由参数
func (this *HttpContext) AddRouterParams(key, value string) {
	this.Params[key] = value
}

// RemoveRouterParams 移除路由参数
func (this *HttpContext) RemoveRouterParams(key string) {
	delete(this.Params, key)
}

// AddRouter 添加执行路由,最后一级路由最先添加
func (this *HttpContext) AddRouter(router router.Router) {
	this.routers = append(this.routers, router)
}

// AddContextExector 添加执行器
func (this *HttpContext) AddContextExecutor(exector router.ContextExecutor) {
	this.executor = exector
}

// ParseParams 解析参数,将路由参数,query string,表单都解析到this.Request.Form中
func (this *HttpContext) ParseParams() error {
	if !this.parsed {
		this.parsed = true
		var err = this.request.ParseForm()
		if err != nil {
			return err
		}
		for k, v := range this.Params {
			this.request.Form.Set(k, v)
		}
	}
	return nil
}

// 处理该HttpContext
func (this *HttpContext) execute() {
	var ok = this.executeBeforeFilters()
	if ok {
		this.executor.Exec(this)
		this.executeAfterFilters()
	}

}

func (this *HttpContext) executeBeforeFilters() bool {
	for i := len(this.routers) - 1; i >= 0; i-- {
		var router = this.routers[i]
		var ok = router.ExecBeforeFilter(this)
		if !ok {
			return false
		}
	}
	return true
}

func (this *HttpContext) executeAfterFilters() bool {
	for i := 0; i < len(this.routers); i++ {
		var router = this.routers[i]
		var ok = router.ExecAfterFilter(this)
		if !ok {
			return false
		}
	}
	return true
}
