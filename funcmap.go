package tinygo

import (
	"html/template"
)

// 模版方法映射接口
type UserFuncMap interface {
	FuncMap() template.FuncMap
}

// 公共模版方法
//  until(num int) []int:生成一个从0到num的数组
type CommonFunMap struct {
}

func (this *CommonFunMap) FuncMap() template.FuncMap {
	return template.FuncMap{
		"until": func(num int) []int {
			var result = make([]int, num)
			for i := 0; i < num; i++ {
				result[i] = i
			}
			return result
		},
	}
}

// CSRF模版方法
//  csrf() string:生成一个csrf字符串
//  csrfHtml() string:生成一个包含csrf的隐藏域
type CsrfFuncMap struct {
	context *HttpContext
}

// Csrf 生成一个csrf字符串
func (this *CsrfFuncMap) Csrf() template.HTML {
	return template.HTML(this.context.CsrfToken())
}

// CsrfHtml 生成一个包含csrf的隐藏域
func (this *CsrfFuncMap) CsrfHtml() template.HTML {
	return template.HTML(`<input type="hidden" name="` + DefaultCSRFTokenName + `" value="` + this.context.CsrfToken() + `" >`)
}

func (this *CsrfFuncMap) FuncMap() template.FuncMap {
	return template.FuncMap{
		"csrf":     this.Csrf,
		"csrfhtml": this.CsrfHtml,
	}
}
