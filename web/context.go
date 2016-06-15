package web

import (
	"io"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"os"
	"reflect"
	"strconv"

	"github.com/kdada/tinygo/connector"
	"github.com/kdada/tinygo/router"
	"github.com/kdada/tinygo/session"
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
	Csrf        session.Session        //csrf会话
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
