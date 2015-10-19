package router

import (
	"net/http"

	"github.com/kdada/tinygo/session"
)

// 路由环境接口
type RouterContext interface {
	// RouterParts 返回路由段
	RouterParts() []string
	// AddRouterParams 添加路由参数
	AddRouterParams(key, value string)
}

// 路由环境
type HttpContext struct {
	UrlParts       []string            //Url分段信息,每段一个字符串
	Request        *http.Request       //http请求
	ResponseWriter http.ResponseWriter //http响应
	Session        session.Session     //http会话
	Params         map[string]string   //http参数,包含url,query,form的所有参数
	parsed         bool                //存储参数是否已经解析过
}

// RouterParts 返回路由段
func (this *HttpContext) RouterParts() []string {
	return this.UrlParts
}

// AddParams 添加路由参数
func (this *HttpContext) AddRouterParams(key, value string) {
	this.Params[key] = value
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
