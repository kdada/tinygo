package tinygo

import (
	"html/template"
	"time"
)

// 模版方法映射接口
type UserFuncMap interface {
	FuncMap() template.FuncMap
}

// 公共模版方法
//  until(num int) []int:生成一个从0到num的数组
//  addi(num ...int) int:整数加法
//  addf(num ...float64) float64:浮点数加法
//  muli(num ...int) int:整数乘法
//  mulf(num ...float64) float64:浮点数乘法
//  tocss(s string) template.CSS:转换字符串为CSS
//  tohtml(s string) template.HTML:转换字符串为HTML
//  toattr(s string) template.HTMLAttr:转换字符串为HTMLAttr
//  tojs(s string) template.JS:转换字符串为JS
//  tojsstr(s string) template.JSStr:转换字符串为JSStr
//  tourl(s string) template.URL:转换字符串为URL
//  time(t time.Time) string:返回时间字符串2006-01-02 15:04:05
//  date(t time.Time) string:返回日期字符串2006-01-02
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
		"addi": func(nums ...int) int {
			var result = 0
			for _, v := range nums {
				result += v
			}
			return result
		},
		"addf": func(nums ...float64) float64 {
			var result = 0.0
			for _, v := range nums {
				result += v
			}
			return result
		},
		"muli": func(nums ...int) int {
			var result = 0
			for _, v := range nums {
				result *= v
			}
			return result
		},
		"mulf": func(nums ...float64) float64 {
			var result = 0.0
			for _, v := range nums {
				result *= v
			}
			return result
		},
		"tocss": func(s string) template.CSS {
			return template.CSS(s)
		},
		"tohtml": func(s string) template.HTML {
			return template.HTML(s)
		},
		"toattr": func(s string) template.HTMLAttr {
			return template.HTMLAttr(s)
		},
		"tojs": func(s string) template.JS {
			return template.JS(s)
		},
		"tojsstr": func(s string) template.JSStr {
			return template.JSStr(s)
		},
		"tourl": func(s string) template.URL {
			return template.URL(s)
		},
		"time": func(t time.Time) string {
			return t.Format("2006-01-02 15:04:05")
		},
		"date": func(t time.Time) string {
			return t.Format("2006-01-02")
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
