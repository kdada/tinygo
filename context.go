package tinygo

import (
	"net/http"

	"github.com/kdada/tinygo/router"
	"github.com/kdada/tinygo/session"
)

// 路由环境
type HttpContext struct {
	UrlParts       []string               //Url分段信息,每段一个字符串
	Request        *http.Request          //http请求
	ResponseWriter http.ResponseWriter    //http响应
	Session        session.Session        //http会话
	Params         map[string]string      //http参数,包含url,query,form的所有参数
	parsed         bool                   //存储参数是否已经解析过
	routers        []router.Router        //分派成功的路由链
	executor       router.ContextExecutor //存储最终执行Context的执行器
}

// RouterParts 返回路由段
func (this *HttpContext) RouterParts() []string {
	return this.UrlParts
}

// SetRouterParts 设置路由段
func (this *HttpContext) SetRouterParts(parts []string) {
	this.UrlParts = parts
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
		var err = this.Request.ParseForm()
		if err != nil {
			return err
		}
		for k, v := range this.Params {
			this.Request.Form.Set(k, v)
		}
	}
	return nil
}
