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
type commonHttpResult struct {
	Status      StatusCode //状态码
	ContentType string     //内容类型
}

// Code 返回状态码
func (this *commonHttpResult) Code() StatusCode {
	return this.Status
}

// Message 返回状态信息
func (this *commonHttpResult) Message() string {
	return ""
}

// WriteHeader 写入响应头信息
func (this *commonHttpResult) WriteHeader(writer io.Writer) (http.ResponseWriter, error) {
	var w, ok = writer.(http.ResponseWriter)
	if !ok {
		return nil, ErrorInvalidWriter.Error()
	}
	if this.ContentType != "" {
		w.Header().Set("Content-Type", this.ContentType)
	}
	w.WriteHeader(int(this.Status))
	return w, nil
}

// 文件结果
type FileResult struct {
	commonHttpResult
	context  *Context //请求上下文
	filePath string   //文件路径
}

// WriteTo 将Result的内容写入writer
func (this *FileResult) WriteTo(writer io.Writer) error {
	var w, ok = writer.(http.ResponseWriter)
	if !ok {
		return ErrorInvalidWriter.Error()
	}
	var f, e = os.Open(this.filePath)
	if e != nil {
		return e
	}
	var info, err = f.Stat()
	if err != nil {
		return err
	}
	if info.IsDir() {
		err = this.dir(this.context.HttpContext.Request, w, f, info)
	} else {
		err = this.file(this.context.HttpContext.Request, w, f, info)
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
	commonHttpResult
	obj interface{} //需要返回的对象
}

// WriteTo 将Result的内容写入writer
func (this *JsonResult) WriteTo(writer io.Writer) error {
	this.ContentType = "application/json; charset=utf-8"
	var w, err = this.commonHttpResult.WriteHeader(writer)
	if err != nil {
		return err
	}
	var bytes, e = json.Marshal(this.obj)
	if e != nil {
		return e
	}
	_, e = w.Write(bytes)
	return e
}

// Xml结果
type XmlResult struct {
	commonHttpResult
	obj interface{} //需要返回的对象
}

// WriteTo 将Result的内容写入writer
func (this *XmlResult) WriteTo(writer io.Writer) error {
	this.ContentType = "application/xml; charset=utf-8"
	var w, err = this.commonHttpResult.WriteHeader(writer)
	if err != nil {
		return err
	}
	var bytes, e = xml.Marshal(this.obj)
	if e != nil {
		return e
	}
	_, e = w.Write(bytes)
	return e
}

// 404结果
type NotFoundResult struct {
	commonHttpResult
	context *Context //请求上下文
}

// WriteTo 将Result的内容写入writer
func (this *NotFoundResult) WriteTo(writer io.Writer) error {
	var r, ok = writer.(http.ResponseWriter)
	if !ok {
		return ErrorInvalidWriter.Error()
	}
	http.NotFound(r, this.context.HttpContext.Request)
	return nil
}

// 重定向结果
type RedirectResult struct {
	commonHttpResult
	context *Context //请求上下文
	url     string   //重定向地址
}

// WriteTo 将Result的内容写入writer
func (this *RedirectResult) WriteTo(writer io.Writer) error {
	var r, ok = writer.(http.ResponseWriter)
	if !ok {
		return ErrorInvalidWriter.Error()
	}
	http.Redirect(r, this.context.HttpContext.Request, this.url, int(this.Status))
	return nil
}

// 数据结果
type DataResult struct {
	commonHttpResult
	data []byte
}

// WriteTo 将Result的内容写入writer
func (this *DataResult) WriteTo(writer io.Writer) error {
	var w, err = this.commonHttpResult.WriteHeader(writer)
	if err != nil {
		return err
	}
	_, err = w.Write(this.data)
	return err
}

// 视图结果
type ViewResult struct {
	commonHttpResult
	templates *ViewTemplates
	path      string
	data      interface{}
}

// WriteTo 将Result的内容写入writer
func (this *ViewResult) WriteTo(writer io.Writer) error {
	var w, err = this.commonHttpResult.WriteHeader(writer)
	if err != nil {
		return err
	}
	return this.templates.ExecView(w, this.path, this.data)
}

// 部分视图结果
type PartialViewResult struct {
	commonHttpResult
	templates *ViewTemplates
	path      string
	data      interface{}
}

// WriteTo 将Result的内容写入writer
func (this *PartialViewResult) WriteTo(writer io.Writer) error {
	var w, err = this.commonHttpResult.WriteHeader(writer)
	if err != nil {
		return err
	}
	return this.templates.ExecPartialView(w, this.path, this.data)
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
