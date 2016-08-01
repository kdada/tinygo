package web

import (
	"html/template"
	"reflect"
	"time"

	"github.com/kdada/tinygo/session"
	"github.com/kdada/tinygo/util"
)

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
var commonFuncMap = template.FuncMap{
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

// RegisterTemplateFunc 注册模板函数,f为nil将删除指定名称的模板函数
func RegisterTemplateFunc(name string, f interface{}) error {
	if f == nil {
		delete(commonFuncMap, name)
		return nil
	}
	var t = reflect.TypeOf(f)
	if t.Kind() == reflect.Func {
		commonFuncMap[name] = f
		return nil
	}
	return ErrorParamMustBeFunc.Error()
}

// 模板会话信息
type TemplateSession struct {
	sess session.Session //当前会话
}

// NewTemplateSession 创建一个模板Session
func NewTemplateSession(sess session.Session) *TemplateSession {
	return &TemplateSession{sess}
}

// String 获取字符串
func (this *TemplateSession) String(key string) (string, error) {
	var v, ok = this.sess.String(key)
	if ok {
		return v, nil
	}
	return "", ErrorInvalidKey.Format(key).Error()
}

//  Int 获取整数值
func (this *TemplateSession) Int(key string) (int, error) {
	var v, ok = this.sess.Int(key)
	if ok {
		return v, nil
	}
	return 0, ErrorInvalidKey.Format(key).Error()
}

// Bool 获取bool值
func (this *TemplateSession) Bool(key string) (bool, error) {
	var v, ok = this.sess.Bool(key)
	if ok {
		return v, nil
	}
	return false, ErrorInvalidKey.Format(key).Error()
}

// Float 获取浮点值
func (this *TemplateSession) Float(key string) (float64, error) {
	var v, ok = this.sess.Float(key)
	if ok {
		return v, nil
	}
	return 0.0, ErrorInvalidKey.Format(key).Error()
}

// 模板CSRF信息
type TemplateCSRF struct {
	sess session.Session //当前CSRF会话
	name string          //token字段名称
}

// NewTemplateCSRF 创建一个模板CSRF
func NewTemplateCSRF(sess session.Session, fieldName string) *TemplateCSRF {
	return &TemplateCSRF{sess, fieldName}
}

// Token 生成一个CSRF认证字符串
func (this *TemplateCSRF) Token() template.HTML {
	var token = util.NewUUID().Hex()
	this.sess.SetInt(token, int(time.Now().Unix())) //记录生成时间(秒)
	return template.HTML(token)
}

// Field 生成一个包含CSRF的隐藏域
func (this *TemplateCSRF) Field() template.HTML {
	return template.HTML(`<input type="hidden" name="` + this.name + `" value="` + string(this.Token()) + `" >`)
}
