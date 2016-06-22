package web

import (
	"io"
	"text/template"
)

// 视图配置
type ViewConfig struct {
	LayoutMap     map[string]string //为布局文件定义别名,map[布局文件别名]布局文件路径
	DefaultLayout string            //默认布局文件别名
	LayoutSpec    map[string]string //特别指定目录或文件对应的布局文件,map[目录或文件路径]布局文件别名
}

// 视图模板信息
type ViewTemplates struct {
	templates   map[string]*template.Template
	config      *ViewConfig //视图配置
	viewPath    string      //视图文件目录
	partialName string      //部分视图的模板名称
	templateExt string      //模板的扩展名
}

// NewViewTemplates 创建视图模板信息
//  path:模板文件的根目录,通常是views目录
//  config:视图布局配置信息
//  partial:部分视图的模板名称
//  ext:模板文件的扩展名(不是以该扩展名结尾的路径都会被加上该名称)
func NewViewTemplates(path string, config *ViewConfig, partial string, ext string) *ViewTemplates {
	return &ViewTemplates{
		make(map[string]*template.Template),
		config,
		path,
		partial,
		ext,
	}
}

// CompileAll 编译所有视图
func (this *ViewTemplates) CompileAll() error {
	return nil
}

// ExecView 执行视图
func (this *ViewTemplates) ExecView(w io.Writer, path string, data interface{}) error {
	return nil
}

// ExecPartialView 执行部分视图
func (this *ViewTemplates) ExecPartialView(w io.Writer, path string, data interface{}) error {
	return nil
}
