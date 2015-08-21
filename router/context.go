package router

import (
	"net/http"
	"tinygo/session"
)

//路由环境接口
type IRouterContext interface {
	// RouterParts 返回路由分段信息
	// 例如 /User/Admin/Test
	// 分段 []string{"User","Admin","Test"}
	RouterParts() []string
}

// 路由环境
type HttpContext struct {
	UrlParts       []string            //Url分段信息,每段一个字符串
	Request        *http.Request       //http请求
	ResponseWriter http.ResponseWriter //http响应
	Session        session.ISession    //http会话
	Params         map[string]string   //http参数,包含url,query,form的所有参数
	parsed         bool                //存储参数是否已经解析过
}

// RouterParts 返回路由分段信息
// 例如 /User/Admin/Test
// 分段 []string{"User","Admin","Test"}
func (this *HttpContext) RouterParts() []string {
	return this.UrlParts
}

// AddParams 添加参数
func (this *HttpContext) AddParams(key, value string) {
	this.Params[key] = value
}

// ParseParams 解析参数,将所有参数解析到this.Request.Form中
func (this *HttpContext) ParseParams() {
	if !this.parsed {
		this.parsed = true
		this.Request.ParseForm()
		for k, v := range this.Params {
			this.Request.Form.Set(k, v)
		}
	}
}
