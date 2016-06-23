package web

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
//  recipi(num int) float64:整数倒数
//  recipf(num float64) float64:浮点数倒数
//  invi(num int) int:整数相反数
//  invf(num float64) float64:浮点数相反数
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
			var result = 1
			for _, v := range nums {
				result *= v
			}
			return result
		},
		"mulf": func(nums ...float64) float64 {
			var result = 1.0
			for _, v := range nums {
				result *= v
			}
			return result
		},
		"recipi": func(num int) float64 {
			if num != 0 {
				return 1.0 / float64(num)
			}
			return 0.0
		},
		"recipf": func(num float64) float64 {
			if num != 0 {
				num = 1.0 / num
			}
			return num
		},
		"invi": func(num int) int {
			return -num
		},
		"invf": func(num float64) float64 {
			return -num
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
