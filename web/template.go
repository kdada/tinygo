package web

import (
	"html/template"
	"time"

	"github.com/kdada/tinygo/session"
	"github.com/kdada/tinygo/util"
)

// 模板会话信息
type TemplateSession struct {
	sess session.Session //当前会话
}

// NewTemplateSession 创建一个模板Session
func NewTemplateSession(sess session.Session) *TemplateSession {
	return &TemplateSession{sess}
}

// Contains 确认Session中是否包含该key
func (this *TemplateSession) Contains(key string) bool {
	var _, ok = this.sess.Value(key)
	return ok
}

// Value 获取值
func (this *TemplateSession) Value(key string) interface{} {
	var v, ok = this.sess.Value(key)
	if ok {
		return v
	}
	return nil
}

// String 获取字符串
func (this *TemplateSession) String(key string) string {
	var v, ok = this.sess.String(key)
	if ok {
		return v
	}
	return ""
}

//  Int 获取整数值
func (this *TemplateSession) Int(key string) int {
	var v, ok = this.sess.Int(key)
	if ok {
		return v
	}
	return 0
}

// Bool 获取bool值
func (this *TemplateSession) Bool(key string) bool {
	var v, ok = this.sess.Bool(key)
	if ok {
		return v
	}
	return false
}

// Float 获取浮点值
func (this *TemplateSession) Float(key string) float64 {
	var v, ok = this.sess.Float(key)
	if ok {
		return v
	}
	return 0.0
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
