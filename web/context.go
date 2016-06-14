package web

import (
	"io"
	"mime/multipart"
	"net/textproto"
	"os"
	"reflect"
	"strconv"

	"github.com/kdada/tinygo/connector"
	"github.com/kdada/tinygo/router"
)

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
	var lf, err = os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0666)
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
	End         router.Router          //路由
	Data        map[string]string      //路由信息
	HttpContext *connector.HttpContext //http上下文
	Processor   *HttpProcessor         //生成当前上下文的处理器
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

// Param 根据名称和类型生成相应类型的数据
func (this *Context) Param(name string, t reflect.Type) interface{} {
	var f, ok = this.Processor.Funcs[t.String()]
	if !ok {
		f = this.Processor.DefaultFunc
	}
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
	var result, ok2 = this.HttpContext.Request.Form[key]
	if ok2 && len(result) > 0 {
		return result[0], nil
	}
	return "", ErrorParamNotExist.Error()
}

// ParamString 获取http参数字符串数组
func (this *Context) ParamStringArray(key string) ([]string, error) {
	var result, ok = this.HttpContext.Request.Form[key]
	var routerResult, ok2 = this.Value(key)
	if !ok && !ok2 {
		return nil, ErrorParamNotExist.Error()
	}
	if !ok {
		result = []string{}
	}
	if ok2 {
		result = append(result, routerResult)
	}
	return result, nil
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

// WriteString 将字符串写入http response流
func (this *Context) WriteString(value string) error {
	var _, err = this.HttpContext.ResponseWriter.Write([]byte(value))
	return err
}
