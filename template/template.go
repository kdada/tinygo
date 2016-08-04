package template

import (
	"html/template"
	"io"
	"os"
	"path/filepath"
)

// 视图模板信息
type ViewTemplates struct {
	templates map[string]*template.Template
	config    *TemplateConfig  //视图配置
	funcMap   template.FuncMap //模板方法
}

// NewViewTemplates 创建视图模板信息
//  config:视图布局配置信息
func NewViewTemplates(config *TemplateConfig) *ViewTemplates {
	return &ViewTemplates{
		make(map[string]*template.Template),
		config,
		commonFuncMap,
	}
}

// CompileAll 编译所有视图
func (this *ViewTemplates) CompileAll() error {
	var templates = make(map[string]*template.Template)
	var err = filepath.Walk(this.config.BasePath(), func(filePath string, fileInfo os.FileInfo, err error) error {
		//遍历目录下的所有模板文件
		if err == nil && fileInfo != nil && !fileInfo.IsDir() && this.config.IsTemplate(fileInfo.Name()) {
			var path = this.config.RelPath(filePath)
			if !this.config.IsLayout(path) {
				var tmpl, err = this.compile(path)
				if err == nil {
					templates[path] = tmpl
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

// Compile 编译指定路径的视图
func (this *ViewTemplates) compile(path string) (*template.Template, error) {
	var pathSlice = this.config.FileLayouts(path)
	var tmplName = filepath.Base(pathSlice[len(pathSlice)-1])
	var tmpl = template.New(tmplName)
	//增加模版方法
	if this.funcMap != nil {
		tmpl.Funcs(this.funcMap)
	}
	return tmpl.ParseFiles(pathSlice...)
}

// template 返回指定路径的视图模板,如果模板不存在则编译该模板
func (this *ViewTemplates) template(path string) (*template.Template, error) {
	path = this.config.Clean(path)
	var tmpl, ok = this.templates[path]
	if ok {
		return tmpl, nil
	}
	return this.compile(path)
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
	content := tmpl.Lookup(this.config.partialName)
	if content != nil {
		return content.Execute(w, data)
	}
	return ErrorInvalidPartialView.Format(path, this.config.partialName).Error()
}
