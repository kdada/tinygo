package tinygo

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"strings"

	"github.com/kdada/tinygo/info"
	"github.com/kdada/tinygo/router"
)

//控制器接口
type IController interface {
	// SetContext 设置请求上下文环境
	SetContext(context *router.HttpContext)
	// SetRouter 设置使用当前控制器的路由
	SetRouter(router router.Router)
	// File 返回文件
	File(path string)
	// Json 返回json格式的数据
	Json(value interface{})
	// Xml 返回Xml格式的数据
	Xml(value interface{})
	// Api 根据设置返回Json或Xml
	Api(value interface{})
	// View 返回视图页面
	View(path string, data ...interface{})
	// SimpleView 返回 控制器名(不含Controller)/方法名.html 页面
	SimpleView(data ...interface{})
	// PartialView 返回 控制器名(不含Controller)/方法名.html 页面无视layout设置
	PartialView(path string, data ...interface{})
	// PartialView 返回 控制器名(不含Controller)/方法名.html 页面无视layout设置
	SimplePartialView(data ...interface{})
	// HttpNotFound 返回404
	HttpNotFound()
	// ParseParams 将参数解析到结构体中
	// params:结构体指针数组,参数必须是结构体指针
	ParseParams(params ...interface{})
	// RedirectMethod [302] 重定向到当前控制器的方法
	// method:方法名
	// params:要传递的参数(这些参数将作为query string传递)
	RedirectMethod(method string, params ...interface{})
	// Redirect [302] 重定向到指定控制器的指定方法
	// controller:控制器名,该控制器必须与当前控制器处于同一个SpaceRouter中
	// method:方法名
	// params:要传递的参数(这些参数将作为query string传递)
	Redirect(controller string, method string, params ...interface{})
	// Redirect [302] 重定向到指定url
	RedirectUrl(url string)
	// RedirectPermanently [301] 永久重定向到指定控制器的方法
	// controller:控制器名,该控制器必须与当前控制器处于同一个SpaceRouter中
	// method:方法名
	// params:要传递的参数(这些参数将作为query string传递)
	RedirectPermanently(controller string, method string, params ...interface{})
	// RedirectUrlPermanently [301] 永久重定向到指定url
	RedirectUrlPermanently(url string)
	// Routers 返回当前控制器可以使用的方法路由信息
	Routers() []RouterInfo
}

//视图数据类型
type ViewData map[interface{}]interface{}

//参数数据类型
type ParamData map[string]string

// 控制器
type Controller struct {
	Context *router.HttpContext //环境
	Router  router.Router       //选择了当前控制器的最后一级路由
	Data    ViewData            //用于传递给页面的数据
}

// SetContext 设置请求上下文环境
func (this *Controller) SetContext(context *router.HttpContext) {
	this.Context = context
	this.Data = make(map[interface{}]interface{}, 0)
}

// SetRouter 设置使用当前控制器的路由
func (this *Controller) SetRouter(router router.Router) {
	this.Router = router
}

// File 返回文件
func (this *Controller) File(path string) {
	http.ServeFile(this.Context.ResponseWriter, this.Context.Request, path)
}

// Json 返回json格式的数据
func (this *Controller) Json(value interface{}) {
	var bytes, err = json.Marshal(value)
	if err != nil {
		fmt.Println(err)
		this.HttpNotFound()
	} else {
		this.Context.ResponseWriter.Header().Set("Content-Type", "application/json")
		_, err := this.Context.ResponseWriter.Write(bytes)
		if err != nil {
			fmt.Println(err)
			this.Context.ResponseWriter.WriteHeader(404)
		}
	}
}

// Xml 返回Xml格式的数据
// 结构体不得是匿名结构体,Xml解析会出错
func (this *Controller) Xml(value interface{}) {
	var bytes, err = xml.Marshal(value)
	if err != nil {
		fmt.Println(err)
		this.HttpNotFound()
	} else {
		this.Context.ResponseWriter.Header().Set("Content-Type", "application/xml")
		_, err := this.Context.ResponseWriter.Write(bytes)
		if err != nil {
			fmt.Println(err)
			this.Context.ResponseWriter.WriteHeader(404)
		}
	}
}

// Api 根据设置返回Json或Xml
func (this *Controller) Api(value interface{}) {
	var api = info.ApiTypeJson
	if info.ApiType(tinyConfig.api) == info.ApiTypeAuto {
		//检测请求头中是否包含指定的api格式
		//优先检测json格式,如果存在指定格式则返回指定格式
		//如果均不存在则返回json格式
		var accept = this.Context.Request.Header.Get("Accept")
		var posJson = strings.Index(accept, "application/json")
		if posJson > 0 {
			api = info.ApiTypeJson
		} else {
			var posXml = strings.Index(accept, "application/xml")
			if posXml > 0 {
				api = info.ApiTypeXml
			}
		}
	}
	switch api {
	case info.ApiTypeJson:
		{
			this.Json(value)
		}
	case info.ApiTypeXml:
		{
			this.Xml(value)
		}
	default:
		{

		}
	}
}

// SetData 设置数据到this.Data中
func (this *Controller) SetData(data ...interface{}) {
	for _, output := range data {
		var outputType = reflect.TypeOf(output)
		switch {
		case outputType.Kind() == reflect.Map:
			{
				//添加Map
				var dataMap, ok = output.(ViewData)
				if ok {
					for k, v := range dataMap {
						this.Data[k] = v
					}
				}
			}
		case outputType.Kind() == reflect.Ptr && outputType.Elem().Kind() == reflect.Struct:
			{
				//将结构体的字段反射到this.Data中
				mapStructToMap(reflect.ValueOf(output).Elem(), this.Data)
			}
		}
	}
}

// DefaultViewPath 返回当前默认的视图路径 控制器名(不含Controller)/方法名.html
func (this *Controller) DefaultViewPath() string {
	var filePath = this.Router.Super().Name() + "/" + this.Router.Name()
	if !strings.HasSuffix(filePath, info.DefaultTemplateExt) {
		filePath += info.DefaultTemplateExt
	}
	return filePath
}

// View 返回视图页面
// path:网页文件相对tinygo.ViewPath目录的位置,例 admin/login.html
// data:需要传递给网页的结构体(必须是指针)或map(必须是ViewData类型),能传递的字段必须是公开字段
func (this *Controller) View(path string, data ...interface{}) {
	this.SetData(data...)
	ParseTemplate(this.Context.ResponseWriter, this.Context.Request, path, this.Data)
}

// SimpleView 返回 控制器名(不含Controller)/方法名.html 页面
func (this *Controller) SimpleView(data ...interface{}) {
	this.View(this.DefaultViewPath(), data...)
}

// PartialView 返回 控制器名(不含Controller)/方法名.html 页面无视layout设置
func (this *Controller) PartialView(path string, data ...interface{}) {
	this.SetData(data...)
	ParsePartialTemplate(this.Context.ResponseWriter, this.Context.Request, path, this.Data)

}

// PartialView 返回 控制器名(不含Controller)/方法名.html 页面无视layout设置
func (this *Controller) SimplePartialView(data ...interface{}) {
	this.PartialView(this.DefaultViewPath(), data...)
}

// HttpNotFound 返回404
func (this *Controller) HttpNotFound() {
	HttpNotFound(this.Context.ResponseWriter, this.Context.Request)
}

// ParseParams 将参数解析到结构体中
// params:结构体指针数组,参数必须是结构体指针
func (this *Controller) ParseParams(params ...interface{}) {
	for _, param := range params {
		var paramType = reflect.TypeOf(param)
		var paramValue = reflect.ValueOf(param)
		if paramType.Kind() == reflect.Ptr && paramType.Elem().Kind() == reflect.Struct && paramValue.Elem().CanSet() {
			var err = this.Context.Request.ParseForm()
			if err == nil {
				ParseUrlValueToStruct(this.Context.Request.Form, paramValue.Elem())
			} else {
				fmt.Println(err)
			}
		}
	}
}

// DataToUrlParam 将结构体或map转换成query字符串
// data:结构体指针或map[string]string
func (this *Controller) DataToUrlParam(data ...interface{}) string {
	var value url.Values
	for _, output := range data {
		var paramData, ok = output.(ParamData)
		if ok {
			for k, v := range paramData {
				value.Set(k, v)
			}
		}
	}
	return value.Encode()
}

// RedirectMethod [302] 重定向到当前控制器的方法
// method:方法名
// params:要传递的参数(这些参数将作为query string传递)
func (this *Controller) RedirectMethod(method string, params ...interface{}) {
	var url = "/" + method
	var router = this.Router.Super()
	for router.Super() != nil {
		url = "/" + router.Name() + url
		router = router.Super()
	}
	var param = this.DataToUrlParam(params...)
	if param != "" {
		url += "?" + param
	}
	this.RedirectUrl(url)
}

// Redirect [302] 重定向到指定控制器的指定方法
// controller:控制器名,该控制器必须与当前控制器处于同一个SpaceRouter中
// method:方法名
// params:要传递的参数(这些参数将作为query string传递)
func (this *Controller) Redirect(controller string, method string, params ...interface{}) {
	var url = "/" + this.Router.Super().Name() + "/" + method
	var router = this.Router.Super().Super()
	for router.Super() != nil {
		url = "/" + router.Name() + url
		router = router.Super()
	}
	var param = this.DataToUrlParam(params...)
	if param != "" {
		url += "?" + param
	}
	this.RedirectUrl(url)
}

// Redirect [302] 重定向到指定url
func (this *Controller) RedirectUrl(url string) {
	Redirect(this.Context.ResponseWriter, this.Context.Request, url)
}

// RedirectPermanently [301] 永久重定向到指定控制器的方法
// controller:控制器名,该控制器必须与当前控制器处于同一个SpaceRouter中
// method:方法名
// params:要传递的参数(这些参数将作为query string传递)
func (this *Controller) RedirectPermanently(controller string, method string, params ...interface{}) {
	var url = "/" + this.Router.Super().Name() + "/" + method
	var router = this.Router.Super().Super()
	for router.Super() != nil {
		url = "/" + router.Name() + url
		router = router.Super()
	}
	var param = this.DataToUrlParam(params...)
	if param != "" {
		url += "?" + param
	}
	this.RedirectUrlPermanently(url)
}

// RedirectUrlPermanently [301] 永久重定向到指定url
func (this *Controller) RedirectUrlPermanently(url string) {
	RedirectPermanently(this.Context.ResponseWriter, this.Context.Request, url)
}

// Routers 返回当前控制器可以使用的方法路由信息
func (this *Controller) Routers() []RouterInfo {
	return []RouterInfo{}
}
