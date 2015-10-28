package tinygo

import (
	"encoding/json"
	"encoding/xml"
	"net/http"
	"net/url"
	"reflect"
	"strings"

	"github.com/kdada/tinygo/router"
)

// 视图数据类型
type ViewData map[interface{}]interface{}

// 参数数据类型
type ParamData map[string]string

// 基础控制器
type baseController struct {
	Context *HttpContext  //环境
	Router  router.Router //选择了当前控制器的最后一级路由
	Data    ViewData      //用于传递给页面的数据
}

// SetContext 设置请求上下文环境
func (this *baseController) SetContext(context router.RouterContext) {
	var ctx, ok = context.(*HttpContext)
	if ok {
		this.Context = ctx
	} else {
		Error("context严重错误", context)
	}
	this.Data = make(map[interface{}]interface{}, 0)
}

// SetRouter 设置使用当前控制器的路由
func (this *baseController) SetRouter(router router.Router) {
	this.Router = router
}

// File 返回文件
func (this *baseController) File(path string) {
	http.ServeFile(this.Context.responseWriter, this.Context.request, path)
}

// Json 返回json格式的数据
func (this *baseController) Json(value interface{}) {
	var bytes, err = json.Marshal(value)
	if err != nil {
		Error(err)
		this.HttpNotFound()
	} else {
		this.Context.responseWriter.Header().Set("Content-Type", "application/json")
		_, err := this.Context.responseWriter.Write(bytes)
		if err != nil {
			Error(err)
			this.Context.responseWriter.WriteHeader(404)
		}
	}
}

// Xml 返回Xml格式的数据
// 结构体不得是匿名结构体,Xml解析会出错
func (this *baseController) Xml(value interface{}) {
	var bytes, err = xml.Marshal(value)
	if err != nil {
		Error(err)
		this.HttpNotFound()
	} else {
		this.Context.responseWriter.Header().Set("Content-Type", "application/xml")
		_, err := this.Context.responseWriter.Write(bytes)
		if err != nil {
			Error(err)
			this.Context.responseWriter.WriteHeader(404)
		}
	}
}

// Api 根据设置返回Json或Xml
func (this *baseController) Api(value interface{}) {
	var api = ApiTypeJson
	if ApiType(tinyConfig.api) == ApiTypeAuto {
		//检测请求头中是否包含指定的api格式
		//优先检测json格式,如果存在指定格式则返回指定格式
		//如果均不存在则返回json格式
		var accept = this.Context.request.Header.Get("Accept")
		var posJson = strings.Index(accept, "application/json")
		if posJson > 0 {
			api = ApiTypeJson
		} else {
			var posXml = strings.Index(accept, "application/xml")
			if posXml > 0 {
				api = ApiTypeXml
			}
		}
	}
	switch api {
	case ApiTypeJson:
		{
			this.Json(value)
		}
	case ApiTypeXml:
		{
			this.Xml(value)
		}
	default:
		{

		}
	}
}

// SetData 设置数据到this.Data中
func (this *baseController) SetData(data ...interface{}) {
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

// View 返回视图页面
//  path:网页文件相对tinygo.ViewPath目录的位置,例 admin/login.html
//  data:需要传递给网页的结构体(必须是指针)或map(必须是ViewData类型),能传递的字段必须是公开字段
func (this *baseController) View(path string, data ...interface{}) {
	this.SetData(data...)
	ParseTemplate(this.Context, path, this.Data)
}

// PartialView 返回 控制器名(不含Controller)/方法名.html 页面无视layout设置
func (this *baseController) PartialView(path string, data ...interface{}) {
	this.SetData(data...)
	ParsePartialTemplate(this.Context, path, this.Data)

}

// HttpNotFound 返回404
func (this *baseController) HttpNotFound() {
	HttpNotFound(this.Context.responseWriter, this.Context.request)
}

// ParseParams 将参数解析到结构体中
//  params:结构体指针数组,参数必须是结构体指针
func (this *baseController) ParseParams(params ...interface{}) {
	for _, param := range params {
		var paramType = reflect.TypeOf(param)
		var paramValue = reflect.ValueOf(param)
		if paramType.Kind() == reflect.Ptr && paramType.Elem().Kind() == reflect.Struct && paramValue.Elem().CanSet() {
			var err = this.Context.request.ParseForm()
			if err == nil {
				ParseUrlValueToStruct(this.Context.request.Form, paramValue.Elem())
			} else {
				Error(err)
			}
		}
	}
}

// DataToUrlParam 将结构体或map转换成query字符串
//  data:结构体指针或map[string]string
func (this *baseController) DataToUrlParam(data ...interface{}) string {
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

// Redirect [302] 重定向到指定url
func (this *baseController) RedirectUrl(url string) {
	Redirect(this.Context.responseWriter, this.Context.request, url)
}

// RedirectUrlPermanently [301] 永久重定向到指定url
func (this *baseController) RedirectUrlPermanently(url string) {
	RedirectPermanently(this.Context.responseWriter, this.Context.request, url)
}

// 常规控制器
type Controller struct {
	baseController
}

// DefaultViewPath 返回当前默认的视图路径 控制器名(不含Controller)/方法名.html
func (this *Controller) DefaultViewPath() string {
	var filePath = this.Router.Super().Name() + "/" + this.Router.Name()
	if !strings.HasSuffix(filePath, DefaultTemplateExt) {
		filePath += DefaultTemplateExt
	}
	return filePath
}

// SimpleView 返回 控制器名(不含Controller)/方法名.html 页面
func (this *Controller) SimpleView(data ...interface{}) {
	this.View(this.DefaultViewPath(), data...)
}

// PartialView 返回 控制器名(不含Controller)/方法名.html 页面无视layout设置
func (this *Controller) SimplePartialView(data ...interface{}) {
	this.PartialView(this.DefaultViewPath(), data...)
}

// RedirectMethod [302] 重定向到当前控制器的方法
//  method:方法名
//  params:要传递的参数(这些参数将作为query string传递)
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
//  controller:控制器名,该控制器必须与当前控制器处于同一个SpaceRouter中
//  method:方法名
//  params:要传递的参数(这些参数将作为query string传递)
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

// RedirectPermanently [301] 永久重定向到指定控制器的方法
//  controller:控制器名,该控制器必须与当前控制器处于同一个SpaceRouter中
//  method:方法名
//  params:要传递的参数(这些参数将作为query string传递)
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

// Routers 返回当前控制器可以使用的方法路由信息
func (this *Controller) Routers() []interface{} {
	return nil
}

// Restful控制器
type RestfulController struct {
	baseController
}

// Get HTTP GET对应方法
func (this *RestfulController) Get() {

}

// Post HTTP POST对应方法
func (this *RestfulController) Post() {

}

// Put HTTP PUT对应方法
func (this *RestfulController) Put() {

}

// Delete HTTP DELETE对应方法
func (this *RestfulController) Delete() {

}
