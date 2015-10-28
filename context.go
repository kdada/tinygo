package tinygo

import (
	"mime/multipart"
	"net/http"
	"strconv"
	"time"

	"github.com/kdada/tinygo/router"
	"github.com/kdada/tinygo/session"
)

// 路由环境
type HttpContext struct {
	urlParts       []string               //Url分段信息,每段一个字符串
	request        *http.Request          //http请求
	responseWriter http.ResponseWriter    //http响应
	session        session.Session        //http会话
	csrf           session.Session        //csrf会话
	params         map[string]string      //http参数,包含url,query,form的所有参数
	parsed         bool                   //存储参数是否已经解析过
	routers        []router.Router        //分派成功的路由链
	executor       router.ContextExecutor //存储最终执行Context的执行器
	static         bool                   //是否是静态路由
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

// CsrfToken生成一个新的csrf token
func (this *HttpContext) CsrfToken() string {
	if this.csrf != nil {
		var newToken = session.Guid()
		this.csrf.SetInt(newToken, time.Now().Unix())
		return newToken
	}
	return ""
}

// ValidateCsrfToken 验证表单请求中是否存在csrf的token并且该token有效,验证后token即失效
func (this *HttpContext) ValidateCsrfToken() bool {
	if this.csrf != nil {
		var token = this.ParamString(DefaultCSRFTokenName)
		if token != "" {
			var t, ok = this.csrf.Int(token)
			if ok {
				this.csrf.Delete(token)
				return time.Now().Unix()-t <= tinyConfig.csrfexpire
			}
		}
	}
	return false
}

// Session 返回Session
func (this *HttpContext) Session() session.Session {
	return this.session
}

// Cookie 返回指定cookie的值
func (this *HttpContext) Cookie(name string) (string, error) {
	var cookie, err = this.request.Cookie(name)
	if err == nil {
		return cookie.Value, nil
	}
	return "", err
}

// AddSimpleCookie 添加一个简单的cookie,该cookie使用/作为path
//  name:最好只使用英文和数字作为名称,不得包含换行符回车符分号冒号等http特殊字符
//  value:最好只使用英文和数字作为值,不得包含换行符回车符分号冒号等http特殊字符
func (this *HttpContext) AddSimpleCookie(name string, value string, age int) {
	var cookie = new(http.Cookie)
	cookie.Name = name
	cookie.Value = value
	if age > 0 {
		cookie.MaxAge = age
		cookie.Expires = time.Now().Add(time.Duration(cookie.MaxAge) * time.Second)
	}
	cookie.HttpOnly = false
	cookie.Path = "/"
	this.responseWriter.Header().Add("Set-Cookie", cookie.String())
}

// AddCookie 添加一个cookie,cookie信息必须完整,若设置了cookie.MaxAge则必须设置cookie.Expires
func (this *HttpContext) AddCookie(cookie *http.Cookie) {
	this.responseWriter.Header().Add("Set-Cookie", cookie.String())
}

// ParseParams 解析参数,将路由参数,query string,表单都解析到this.Request.Form中
func (this *HttpContext) ParseParams() error {
	if !this.parsed {
		this.parsed = true
		var err = this.request.ParseForm()
		if err != nil {
			return err
		}
		for k, v := range this.params {
			this.request.Form.Set(k, v)
		}
	}
	return nil
}

// ParamString 获取http参数字符串
func (this *HttpContext) ParamString(key string) string {
	return this.request.FormValue(key)
}

// ParamString 获取http参数字符串数组
func (this *HttpContext) ParamStringArray(key string) ([]string, error) {
	var result, ok = this.request.Form[key]
	if ok {
		return result, nil
	}
	return nil, TinyGoErrorParamNotFoundError.Format(key).Error()
}

// ParamString 获取http参数字符串
func (this *HttpContext) ParamBool(key string) (bool, error) {
	var result = this.request.FormValue(key)
	return strconv.ParseBool(result)
}

// ParamString 获取http参数字符串
func (this *HttpContext) ParamInt(key string) (int64, error) {
	var result = this.request.FormValue(key)
	return strconv.ParseInt(result, 10, 64)

}

// ParamString 获取http参数字符串
func (this *HttpContext) ParamFloat(key string) (float64, error) {
	var result = this.request.FormValue(key)
	return strconv.ParseFloat(result, 64)

}

// ParamString 获取http参数文件
func (this *HttpContext) ParamFile(key string) (multipart.File, *multipart.FileHeader, error) {
	return this.request.FormFile(key)
}

// WriteString 将字符串写入http response流
func (this *HttpContext) WriteString(value string) error {
	var _, err = this.responseWriter.Write([]byte(value))
	return err
}

/////////以下为路由使用方法,一般不需要使用/////////

// RouterParts 返回路由段
func (this *HttpContext) RouterParts() []string {
	return this.urlParts
}

// SetRouterParts 设置路由段
func (this *HttpContext) SetRouterParts(parts []string) {
	this.urlParts = parts
}

// Static 返回是否是静态路由
func (this *HttpContext) Static() bool {
	return this.static
}

// SetStatic 设置当前上下文为静态路由上下文
func (this *HttpContext) SetStatic(static bool) {
	this.static = static
}

// AddParams 添加路由参数
func (this *HttpContext) AddRouterParams(key, value string) {
	this.params[key] = value
}

// RemoveRouterParams 移除路由参数
func (this *HttpContext) RemoveRouterParams(key string) {
	delete(this.params, key)
}

// AddRouter 添加执行路由,最后一级路由最先添加
func (this *HttpContext) AddRouter(router router.Router) {
	this.routers = append(this.routers, router)
}

// AddContextExector 添加执行器
func (this *HttpContext) AddContextExecutor(exector router.ContextExecutor) {
	this.executor = exector
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
