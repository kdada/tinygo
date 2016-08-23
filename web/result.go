package web

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/kdada/tinygo/template"
)

// 可用于默认http方法(默认为Post)的返回结果
type Result interface {
	// WriteTo 将Result的内容写入writer
	WriteTo(writer io.Writer) error
}

// 可用于Get方法的返回结果
type GetResult Result

// 可用于Post方法的返回结果
type PostResult Result

// 可用于Put方法的返回结果
type PutResult Result

// 可用于Delete方法的返回结果
type DeleteResult Result

// 可用于Options方法的返回结果
type OptionsResult Result

// 可用于Head方法的返回结果
type HeadResult Result

// 可用于Trace方法的返回结果
type TraceResult Result

// 可用于Connect方法的返回结果
type ConnectResult Result

// 可用于Get和Post方法的返回结果
type GetPostResult Result

// 用于http的结果
type HttpResult interface {
	// Code 返回状态码
	Code() StatusCode
	// Message 返回状态信息
	Message() string
	// 实现Result接口
	Result
}

// 公共http结果
type CommonHttpResult struct {
	Status      StatusCode //状态码
	ContentType string     //内容类型
}

// Code 返回状态码
func (this *CommonHttpResult) Code() StatusCode {
	return this.Status
}

// Message 返回状态信息
func (this *CommonHttpResult) Message() string {
	return ""
}

// SetHeader 设置响应头信息
func (this *CommonHttpResult) SetHeader(writer io.Writer) (http.ResponseWriter, error) {
	var w, ok = writer.(http.ResponseWriter)
	if !ok {
		return nil, ErrorInvalidWriter.Error()
	}
	if this.ContentType != "" {
		w.Header().Set("Content-Type", this.ContentType)
	}
	return w, nil
}

// WriteHeader 写入响应头信息,对于
func (this *CommonHttpResult) WriteHeader(writer http.ResponseWriter) {
	writer.WriteHeader(int(this.Status))
}

// 文件结果
type FileResult struct {
	CommonHttpResult
	Context  *Context //请求上下文
	FilePath string   //本地文件路径
}

// WriteTo 将Result的内容写入writer
func (this *FileResult) WriteTo(writer io.Writer) error {
	var w, e = this.SetHeader(writer)
	if e != nil {
		return e
	}
	var f, err = os.Open(this.FilePath)
	if err != nil {
		return err
	}
	var info, err2 = f.Stat()
	if err2 != nil {
		return err2
	}
	if info.IsDir() {
		err = this.dir(this.Context.HttpContext.Request, w, f, info)
	} else {
		err = this.file(this.Context.HttpContext.Request, w, f, info)
	}
	return err
}

// file 处理普通文件类型
func (this *FileResult) file(r *http.Request, w http.ResponseWriter, f *os.File, info os.FileInfo) error {
	http.ServeContent(w, r, f.Name(), info.ModTime(), f)
	return nil
}

// html标记替换器
var htmlReplacer = strings.NewReplacer(
	"&", "&amp;",
	"<", "&lt;",
	">", "&gt;",
	`"`, "&#34;",
	"'", "&#39;",
)

// dir 处理目录类型
func (this *FileResult) dir(r *http.Request, w http.ResponseWriter, f *os.File, info os.FileInfo) error {
	var dirs, err = f.Readdir(0)
	if err != nil {
		return err
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprintf(w, "<pre>\n")
	for _, d := range dirs {
		name := d.Name()
		if d.IsDir() {
			name += "/"
		}
		var u = url.URL{Path: name}

		fmt.Fprintf(w, "<a href=\"%s\">%s</a>\n", filepath.Join(r.URL.Path, u.String()), htmlReplacer.Replace(name))
	}
	fmt.Fprintf(w, "</pre>\n")
	fmt.Fprintf(w, `

	`)
	return nil
}

// Json结果
type JsonResult struct {
	CommonHttpResult
	Obj interface{} //需要返回的对象
}

// WriteTo 将Result的内容写入writer
func (this *JsonResult) WriteTo(writer io.Writer) error {
	var w, err = this.SetHeader(writer)
	if err != nil {
		return err
	}
	var bytes, e = json.Marshal(this.Obj)
	if e != nil {
		return e
	}
	this.WriteHeader(w)
	_, e = w.Write(bytes)
	return e
}

// Xml结果
type XmlResult struct {
	CommonHttpResult
	Obj interface{} //需要返回的对象
}

// WriteTo 将Result的内容写入writer
func (this *XmlResult) WriteTo(writer io.Writer) error {
	var w, err = this.SetHeader(writer)
	if err != nil {
		return err
	}
	var bytes, e = xml.Marshal(this.Obj)
	if e != nil {
		return e
	}
	this.WriteHeader(w)
	_, e = w.Write(bytes)
	return e
}

// 404结果
type NotFoundResult struct {
	CommonHttpResult
	Context *Context //请求上下文
}

// WriteTo 将Result的内容写入writer
func (this *NotFoundResult) WriteTo(writer io.Writer) error {
	var w, err = this.SetHeader(writer)
	if err != nil {
		return err
	}
	http.NotFound(w, this.Context.HttpContext.Request)
	return nil
}

// 重定向结果
type RedirectResult struct {
	CommonHttpResult
	Context *Context //请求上下文
	Url     string   //重定向地址
}

// WriteTo 将Result的内容写入writer
func (this *RedirectResult) WriteTo(writer io.Writer) error {
	var w, err = this.SetHeader(writer)
	if err != nil {
		return err
	}
	http.Redirect(w, this.Context.HttpContext.Request, this.Url, int(this.Status))
	return nil
}

// 数据结果
type DataResult struct {
	CommonHttpResult
	Data []byte
}

// WriteTo 将Result的内容写入writer
func (this *DataResult) WriteTo(writer io.Writer) error {
	var w, err = this.SetHeader(writer)
	if err != nil {
		return err
	}
	this.WriteHeader(w)
	_, err = w.Write(this.Data)
	return err
}

// 视图结果
type ViewResult struct {
	CommonHttpResult
	Templates *template.ViewTemplates
	Path      string
	Data      interface{}
}

// WriteTo 将Result的内容写入writer
func (this *ViewResult) WriteTo(writer io.Writer) error {
	var w, err = this.SetHeader(writer)
	if err != nil {
		return err
	}
	this.WriteHeader(w)
	return this.Templates.ExecView(w, this.Path, this.Data)
}

// 部分视图结果
type PartialViewResult struct {
	CommonHttpResult
	Templates *template.ViewTemplates
	Path      string
	Data      interface{}
}

// WriteTo 将Result的内容写入writer
func (this *PartialViewResult) WriteTo(writer io.Writer) error {
	var w, err = this.SetHeader(writer)
	if err != nil {
		return err
	}
	this.WriteHeader(w)
	return this.Templates.ExecPartialView(w, this.Path, this.Data)
}

// 自定义返回结果
type UserDefinedResult struct {
	Status StatusCode //状态码
	Msg    string     //消息
}

// NewUserDefinedResult 创建自定义的返回结果
func NewUserDefinedResult(code StatusCode, msg string) *UserDefinedResult {
	return &UserDefinedResult{
		code,
		msg,
	}
}

// Code 返回状态码
func (this *UserDefinedResult) Code() StatusCode {
	return this.Status
}

// Message 返回状态信息
func (this *UserDefinedResult) Message() string {
	return this.Msg
}

// WriteTo 将Result的内容写入writer
func (this *UserDefinedResult) WriteTo(writer io.Writer) error {
	var _, err = writer.Write([]byte(this.Msg))
	return err
}
