package template

import (
	"path/filepath"
	"strings"
)

// 视图模板配置
type TemplateConfig struct {
	basePath      string            //模板根目录
	partialName   string            //部分模板名称
	templateExt   string            //模板扩展名
	layouts       map[string]string //布局文件,map[布局文件路径]布局文件别名
	DefaultLayout string            //默认布局文件别名,默认布局仅对非布局文件有效
	LayoutMap     map[string]string //为布局文件定义别名,map[布局文件别名]布局文件路径
	LayoutSpec    map[string]string //特别指定目录或文件对应的布局文件,map[目录或文件路径]布局文件别名
}

// 创建模板配置
func NewTemplateConfig() *TemplateConfig {
	return new(TemplateConfig)
}

// SetBasePath 设置模板根目录,LayoutMap和LayoutSpec中的所有文件路径都相对于该路径
func (this *TemplateConfig) SetBasePath(base string) {
	this.basePath = this.Clean(base)
}

// BasePath 返回模板根目录
func (this *TemplateConfig) BasePath() string {
	return this.basePath
}

// SetPartialName 设置部分模板名称,用于在模板中查找部分模板时使用
func (this *TemplateConfig) SetPartialName(name string) {
	this.partialName = name
}

// PartialName 返回部分模板名称
func (this *TemplateConfig) PartialName() string {
	return this.partialName
}

// SetTemplateExt 设置模板扩展名,只有扩展名为该名称的文件才被作为模板
func (this *TemplateConfig) SetTemplateExt(ext string) {
	this.templateExt = strings.ToLower(ext)
}

// TemplateExt 返回模板扩展名
func (this *TemplateConfig) TemplateExt() string {
	return this.templateExt
}

// IsTemplate 判断path指向的文件是否是模板
func (this *TemplateConfig) IsTemplate(path string) bool {
	return strings.ToLower(filepath.Ext(path)) == this.TemplateExt()
}

// IsLayout 判断path指向的文件是否是布局文件,所有路径为相对于BasePath的路径,第一次使用该方法前需要用Clear方法对路径进行规范化
func (this *TemplateConfig) IsLayout(path string) bool {
	if this.layouts != nil {
		var _, ok = this.layouts[this.Clean(path)]
		return ok
	}
	return false
}

// ParentLayout 返回path的父布局路径,所有路径为相对于BasePath的路径
func (this *TemplateConfig) ParentLayout(path string) (string, bool) {
	path = this.Clean(path)
	var layout, ok = this.LayoutSpec[path]
	if !ok {
		var dir = filepath.Dir(path)
		layout, ok = this.LayoutSpec[dir]
		if !ok && !this.IsLayout(path) {
			ok = true
			layout = this.DefaultLayout
		}
	}
	if !ok || layout == "" {
		return "", false
	}
	layout, ok = this.LayoutMap[layout]
	if path == layout {
		return "", false
	}
	return layout, ok
}

// Layouts 返回包含path以及path的父布局路径数组,所有路径为相对于BasePath的路径
func (this *TemplateConfig) Layouts(path string) []string {
	path = this.Clean(path)
	var result = []string{path}
	var ok = true
	for ok {
		path, ok = this.ParentLayout(path)
		if ok {
			result = append(result, path)
		}
	}
	return result
}

// RelPath 返回真实文件路径相对于BasePath的路径
func (this *TemplateConfig) RelPath(realPath string) string {
	var rpath, err = filepath.Rel(this.BasePath(), realPath)
	if err != nil {
		rpath = realPath
	}
	return this.Clean(rpath)
}

// FilePath 返回path的真实文件路径
func (this *TemplateConfig) FilePath(path string) string {
	return filepath.Join(this.BasePath(), path)
}

// Layouts 返回包含path以及path的父布局路径数组,path为相对于BasePath的路径,返回路径为真实文件路径
func (this *TemplateConfig) FileLayouts(path string) []string {
	var ps = this.Layouts(path)
	for i, p := range ps {
		ps[i] = this.FilePath(p)
	}
	return ps
}

// Clean 返回path规范化后的值
func (this *TemplateConfig) Clean(path string) string {
	return strings.Trim(filepath.Clean(path), "/\\")
}

// Clear 将LayoutMap和LayoutSpec中的所有文件路径都进行规范化
func (this *TemplateConfig) Clear() {
	if this.LayoutMap != nil {
		this.layouts = make(map[string]string)
		for k, v := range this.LayoutMap {
			v = this.Clean(v)
			this.LayoutMap[k] = v
			this.layouts[v] = k
		}
	}
	if this.LayoutSpec != nil {
		var spec = make(map[string]string)
		for k, v := range this.LayoutSpec {
			spec[this.Clean(k)] = v
		}
		this.LayoutSpec = spec
	}
}
