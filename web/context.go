package web

import (
	"io"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/kdada/tinygo/connector"
	"github.com/kdada/tinygo/router"
	"github.com/kdada/tinygo/session"
)

// 视图数据类型
type ViewData map[string]interface{}

// 表单文件
type FormFile struct {
	file   multipart.File        //表单文件
	header *multipart.FileHeader //文件头信息
}

// FileName 返回文件名
func (this *FormFile) FileName() string {
	return this.header.Filename
}

// Header 返回文件头信息
func (this *FormFile) Header() textproto.MIMEHeader {
	return this.header.Header
}

// File 返回表单文件
func (this *FormFile) File() multipart.File {
	return this.file
}

// Close 关闭表单文件
func (this *FormFile) Close() {
	this.file.Close()
}

// SaveTo 将文件保存到指定路径,保存完毕后自动关闭表单文件
func (this *FormFile) SaveTo(path string) error {
	var lf, err = os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0660)
	if err == nil {
		defer lf.Close()
		_, err = io.Copy(lf, this.file)
		if err == nil {
			this.Close()

		}
	}
	return err
}

type Context struct {
	router.BaseContext
	HttpContext *connector.HttpContext //http上下文
	Session     session.Session        //http会话
	CSRF        session.Session        //csrf会话
	End         router.Router          //处理当前上下文的路由
	Processor   *HttpProcessor         //生成当前上下文的处理器
}

// NewContext 创建上下文信息
func NewContext(segments []string, context *connector.HttpContext) (*Context, error) {
	var err = context.Request.ParseForm()
	if err != nil {
		return nil, err
	}
	var method = context.Request.Method
	var c = new(Context)
	c.Segs = append(segments, method)
	c.HttpContext = context
	return c, nil
}

// Value 返回值
func (this *Context) Value(name string) (string, bool) {
	return this.HttpContext.Request.Form.Get(name), true
}

// Value 返回值数组
func (this *Context) Values(name string) ([]string, bool) {
	var result, ok = this.HttpContext.Request.Form[name]
	return result, ok
}

// SetValue 设置值
func (this *Context) SetValue(name string, value string) {
	var _, ok = this.HttpContext.Request.Form[name]
	if ok {
		this.HttpContext.Request.Form.Add(name, value)
	} else {
		this.HttpContext.Request.Form.Set(name, value)
	}
}

// Param 根据名称和类型生成相应类型的数据,使用HttpProcessor中定义的参数生成方法
//  name:名称,根据该名称从Request里取值
//  t:生成类型,将对应值转换为该类型
//  return:返回指定类型的数据
func (this *Context) Param(name string, t reflect.Type) interface{} {
	var f = this.Processor.ParamFunc(t.String())
	if f != nil {
		return f(this, name, t)
	}
	return nil
}

// ParamString 获取http参数字符串
func (this *Context) ParamString(key string) (string, error) {
	var routerResult, ok = this.Value(key)
	if ok {
		return routerResult, nil
	}
	return "", ErrorParamNotExist.Error()
}

// ParamString 获取http参数字符串数组
func (this *Context) ParamStringArray(key string) ([]string, error) {
	var result, ok = this.Values(key)
	if ok {
		return result, nil
	}
	return nil, ErrorParamNotExist.Error()
}

// ParamString 获取http参数字符串
func (this *Context) ParamBool(key string) (bool, error) {
	var result, err = this.ParamString(key)
	if err == nil {
		return strconv.ParseBool(result)
	}
	return false, err
}

// ParamString 获取http参数字符串
func (this *Context) ParamInt(key string) (int, error) {
	var result, err = this.ParamString(key)
	if err == nil {
		var v, err2 = strconv.ParseInt(result, 10, 64)
		return int(v), err2
	}
	return 0, err

}

// ParamString 获取http参数字符串
func (this *Context) ParamFloat(key string) (float64, error) {
	var result, err = this.ParamString(key)
	if err == nil {
		return strconv.ParseFloat(result, 64)
	}
	return 0, err
}

// ParamString 获取http参数文件
func (this *Context) ParamFile(key string) (*FormFile, error) {
	var file, header, err = this.HttpContext.Request.FormFile(key)
	if err == nil {
		return &FormFile{file, header}, nil
	}
	return nil, err
}

// ParamString 获取http参数文件数组
func (this *Context) ParamFiles(key string) ([]*FormFile, error) {
	var r = this.HttpContext.Request
	if r.MultipartForm == nil {
		err := r.ParseMultipartForm(int64(this.Processor.Config.MaxRequestMemory))
		if err != nil {
			return nil, err
		}
	}
	if r.MultipartForm != nil && r.MultipartForm.File != nil {
		if fhs := r.MultipartForm.File[key]; len(fhs) > 0 {
			var files = make([]*FormFile, 0)
			for _, fh := range fhs {
				f, err := fh.Open()
				if err == nil {
					files = append(files, &FormFile{f, fh})
				} else {
					return nil, err
				}
			}
			return files, nil
		}
	}
	return nil, http.ErrMissingFile
}

// WriteString 将字符串写入http response流
func (this *Context) WriteString(value string) error {
	var _, err = this.HttpContext.ResponseWriter.Write([]byte(value))
	return err
}

// WriteResult 将Result写入http response流
func (this *Context) WriteResult(result Result) error {
	return result.WriteTo(this.HttpContext.ResponseWriter)
}

// 返回文件类型结果
func (this *Context) File(path string) *FileResult {
	var result = new(FileResult)
	result.Status = 200
	result.ContentType = ContentType(filepath.Ext(path))
	result.context = this
	result.filePath = path
	return result
}

// 返回Json类型结果
func (this *Context) Json(data interface{}) *JsonResult {
	var result = new(JsonResult)
	result.Status = 200
	result.ContentType = ""
	result.obj = data
	return result
}

// 返回Xml类型结果
func (this *Context) Xml(data interface{}) *XmlResult {
	var result = new(XmlResult)
	result.Status = 200
	result.ContentType = ""
	result.obj = data
	return result
}

// 返回Api类型结果
func (this *Context) Api(data interface{}) HttpResult {
	switch this.Processor.Config.Api {
	case "json":
		{
			return this.Json(data)
		}
	case "xml":
		{
			return this.Xml(data)
		}
	default:
		{
			//检测请求头中是否包含指定的api格式
			//优先检测json格式,如果存在指定格式则返回指定格式
			//如果均不存在则返回json格式
			var accept = this.HttpContext.Request.Header.Get("Accept")
			if strings.Index(accept, "application/json") >= 0 {
				return this.Json(data)
			} else if strings.Index(accept, "application/xml") >= 0 {
				return this.Json(data)
			}
		}
	}
	return this.Json(data)
}

// 返回NotFound类型结果
func (this *Context) NotFound() *NotFoundResult {
	var result = new(NotFoundResult)
	result.Status = 404
	result.ContentType = ""
	result.context = this
	return result
}

// 返回临时重定向类型结果
func (this *Context) Redirect(url string) *RedirectResult {
	var result = new(RedirectResult)
	result.Status = 302
	result.ContentType = ""
	result.context = this
	result.url = url
	return result
}

// 返回永久重定向类型结果
func (this *Context) RedirectPermanently(url string) *RedirectResult {
	var result = new(RedirectResult)
	result.Status = 301
	result.ContentType = ""
	result.context = this
	result.url = url
	return result
}

// 返回数据类型结果
func (this *Context) Data(data []byte) *DataResult {
	var result = new(DataResult)
	result.Status = 200
	result.ContentType = ""
	result.data = data
	return result
}

// 返回视图类型结果
func (this *Context) View(path string, data ...interface{}) *ViewResult {
	var result = new(ViewResult)
	result.Status = 200
	result.ContentType = ""
	result.templates = this.Processor.Templates
	result.path = path
	data = append(data, this.commonViewData())
	result.data = this.integrate(data...)
	return result
}

// 返回部分视图类型结果
func (this *Context) PartialView(path string, data ...interface{}) *PartialViewResult {
	var result = new(PartialViewResult)
	result.Status = 200
	result.ContentType = ""
	result.templates = this.Processor.Templates
	result.path = path
	data = append(data, this.commonViewData())
	result.data = this.integrate(data...)
	return result
}

// commonViewData 生成公共视图数据
func (this *Context) commonViewData() ViewData {
	var vd = ViewData{}
	if this.Session != nil {
		vd["SESSION"] = NewTemplateSession(this.Session)
	}
	if this.CSRF != nil {
		vd["CSRF"] = NewTemplateCSRF(this.CSRF, this.Processor.Config.CSRFTokenName)
	}
	return vd
}

// integrate 将data的数据整合为一个ViewData
func (this *Context) integrate(data ...interface{}) ViewData {
	var result = make(ViewData)
	for _, output := range data {
		var outputType = reflect.TypeOf(output)
		if outputType.Kind() == reflect.Map {
			//添加Map
			var dataMap, ok = output.(ViewData)
			if ok {
				for k, v := range dataMap {
					result[k] = v
				}
			}
		} else if IsStructPtrType(outputType) {
			this.mapTo(reflect.ValueOf(output).Elem(), result)
		}
	}
	return result
}

// mapTo 将value结构体转换到一个ViewData中,value必须是结构体
func (this *Context) mapTo(value reflect.Value, data ViewData) {
	if value.Kind() == reflect.Struct {
		for i := 0; i < value.NumField(); i++ {
			var fieldValue = value.Field(i)
			if fieldValue.CanInterface() {
				var fieldType = value.Type().Field(i)
				if fieldType.Anonymous {
					//匿名组合字段,进行递归解析
					this.mapTo(fieldValue, data)
				} else {
					//非匿名字段
					data[fieldType.Name] = fieldValue.Interface()
				}
			}
		}
	}
}

// 重新分发,将当前请求交给另一个path处理
func (this *Context) Redispatch(path string) *UserDefinedResult {
	var result = new(UserDefinedResult)
	result.Status = StatusCodeRedispatch
	result.Msg = path
	return result
}

// ValidateCSRF 验证表单请求中是否存在csrf的token并且该token有效,验证后token立即失效
func (this *Context) ValidateCSRF() bool {
	if this.CSRF != nil {
		var token, err = this.ParamString(this.Processor.Config.CSRFTokenName)
		if err == nil {
			var t, ok = this.CSRF.Int(token)
			if ok {
				this.CSRF.Delete(token)
				return int(time.Now().Unix())-t <= this.Processor.Config.CSRFExpire
			}
		}
	}
	return false
}
