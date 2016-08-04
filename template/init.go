package template

import (
	"html/template"
	"time"
)

func init() {
	//初始化基本模板函数
	//until(num int) []int:生成一个从0到num的数组
	//time(t time.Time) string:返回时间字符串2006-01-02 15:04:05
	//date(t time.Time) string:返回日期字符串2006-01-02
	//tocss(s string) template.CSS:转换字符串为CSS
	//tohtml(s string) template.HTML:转换字符串为HTML
	//toattr(s string) template.HTMLAttr:转换字符串为HTMLAttr
	//tojs(s string) template.JS:转换字符串为JS
	//tojsstr(s string) template.JSStr:转换字符串为JSStr
	//tourl(s string) template.URL:转换字符串为URL
	commonFuncMap = template.FuncMap{
		"until": func(num int) []int {
			var result = make([]int, num)
			for i := 0; i < num; i++ {
				result[i] = i
			}
			return result
		},
		"time": func(t time.Time) string {
			return t.Format("2006-01-02 15:04:05")
		},
		"date": func(t time.Time) string {
			return t.Format("2006-01-02")
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
	}
}
