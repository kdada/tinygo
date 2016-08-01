package web

import (
	"html/template"
	"io"
	"os"
	"path/filepath"
	"strings"
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
	config      *ViewConfig      //视图配置
	viewPath    string           //视图文件目录
	partialName string           //部分视图的模板名称
	templateExt string           //模板的扩展名(小写)
	funcMap     template.FuncMap //最后一次编译使用的模板方法信息
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
		strings.ToLower(ext),
		nil,
	}
}

// CompileAll 编译所有视图
func (this *ViewTemplates) CompileAll(funcMap template.FuncMap) error {
	this.funcMap = funcMap
	var templates = make(map[string]*template.Template)
	var err = filepath.Walk(this.viewPath, func(filePath string, fileInfo os.FileInfo, err error) error {
		//遍历目录下的所有扩展名为templateExt的文件
		if err == nil && fileInfo != nil && !fileInfo.IsDir() && strings.ToLower(filepath.Ext(fileInfo.Name())) == this.templateExt {
			filePath = this.rel(filePath)
			if !this.isLayout(filePath) {
				var tmpl, err = this.compile(filePath, this.funcMap)
				if err == nil {
					templates[filePath] = tmpl
				} else {
					return err
				}
			}
		}
		return nil
	})
	if err == nil {
		this.templates = templates
	}
	return err
}

// rel 返回path相对于viewPath的路径
func (this *ViewTemplates) rel(path string) string {
	var rpath, err = filepath.Rel(this.viewPath, path)
	if err != nil {
		rpath = path
	}
	return NormalizePath(rpath)
}

// layout 获取指定文件的布局文件
func (this *ViewTemplates) layout(path string) (string, bool) {
	var layout, ok = this.config.LayoutSpec[path]
	if !ok {
		var dir = filepath.Dir(path)
		layout, ok = this.config.LayoutSpec[dir]
		if !ok {
			layout = this.config.DefaultLayout
		}
	}
	layout, ok = this.config.LayoutMap[layout]
	if path == layout {
		//避免循环布局
		return "", false
	}
	return layout, ok
}

// isLayout 判断是否是布局文件
func (this *ViewTemplates) isLayout(path string) bool {
	for _, layoutPath := range this.config.LayoutMap {
		if path == layoutPath {
			return true
		}
	}
	return false
}

// file 获取指定视图路径的文件路径
func (this *ViewTemplates) file(path string) string {
	return filepath.Join(this.viewPath, path)
}

// Compile 编译指定路径的视图
func (this *ViewTemplates) compile(path string, funcMap template.FuncMap) (*template.Template, error) {
	var pathSlice = []string{path}
	var layout = path
	var ok = true
	//查找布局文件
	for ok {
		layout, ok = this.layout(layout)
		if ok {
			pathSlice = append(pathSlice, this.file(layout))
		}
	}
	var tmplName = filepath.Base(layout)
	var tmpl = template.New(tmplName)
	//增加模版方法
	tmpl.Funcs(funcMap)
	var tmpls, err = tmpl.ParseFiles(pathSlice...)
	if err == nil {
		var name = filepath.Base(pathSlice[len(pathSlice)-1])
		tmpl = tmpls.Lookup(name)
	}
	return tmpl, err
}

// template 返回指定路径的视图模板,如果模板不存在则编译该模板
func (this *ViewTemplates) template(path string) (*template.Template, error) {
	path = NormalizePath(path)
	var tmpl, ok = this.templates[path]
	if ok {
		return tmpl, nil
	}
	path = this.file(path)
	return this.compile(path, this.funcMap)
}

// ExecView 执行视图
func (this *ViewTemplates) ExecView(w io.Writer, path string, data interface{}) error {
	var tmpl, err = this.template(path)
	if err != nil {
		return err
	}
	return tmpl.Execute(w, data)
}

// ExecPartialView 执行部分视图
func (this *ViewTemplates) ExecPartialView(w io.Writer, path string, data interface{}) error {
	var tmpl, err = this.template(path)
	if err != nil {
		return err
	}
	content := tmpl.Lookup(this.partialName)
	if content != nil {
		return content.Execute(w, data)
	}
	return ErrorInvalidPartialView.Format(path, this.partialName).Error()
}
