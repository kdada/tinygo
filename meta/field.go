package meta

import (
	"reflect"
	"regexp"
	"strings"

	"github.com/kdada/tinygo/validator"
)

// 字段验证类型,字段没有设置验证类型时,等价于没有验证器的可选验证
type FieldKind byte

const (
	FieldKindIgnore   FieldKind = iota //忽略字段 -  ,忽略该字段,不验证也不注入,使用字段初始值
	FieldKindOptional                  //可选验证 ?  ,任何情况下不报告错误,如果有验证器则进行验证,验证通过后注入值,如果没有验证器则直接注入值
	FieldKindRequired                  //必须验证 !  ,未成功注入值则报告错误,如果有验证器则进行验证,验证通过后注入值,如果没有验证器则直接注入值
)

// 字段元数据
type FieldMetadata struct {
	Name      string               //字段名
	Field     *reflect.StructField //字段信息
	Kind      FieldKind            //字段解析类型
	Validator validator.Validator  //验证器
}

// Set 使用vc设置object的值
//  instance:拥有当前字段的对象的指针的反射值
func (this *FieldMetadata) Set(instance reflect.Value, vc ValueContainer) error {
	if this.Kind != FieldKindIgnore {
		var vp, exist = vc.Contains(this.Name, this.Field.Type)
		if exist {
			var valid = true
			if this.Validator != nil {
				//验证器校验
				var strs = vp.String()
				if len(strs) <= 0 && this.Kind == FieldKindRequired {
					return ErrorRequiredField.Format(this.Name).Error()
				}
				for i, v := range strs {
					valid = this.Validator.Validate(v)
					if !valid {
						if this.Kind == FieldKindRequired {
							if len(strs) == 1 {
								return ErrorFieldNotValid.Format(this.Name).Error()
							}
							return ErrorFieldsNotValid.Format(this.Name, i).Error()
						}
						break
					}
				}
			}
			if valid {
				//验证通过则注入值
				var value = vp.Value()
				if value != nil {
					var fValue = reflect.ValueOf(value)
					instance.Elem().FieldByIndex(this.Field.Index).Set(fValue)
				} else if this.Kind == FieldKindRequired {
					return ErrorRequiredField.Format(this.Name).Error()
				}
			}
		} else if this.Kind == FieldKindRequired {
			return ErrorRequiredField.Format(this.Name).Error()
		}
	}
	return nil
}

// 验证字符串提取正则
var vldReg = regexp.MustCompile("^[?!] *?;(.*)$")

// AnalyzeField 分析字段
//  field:结构体字段信息,使用字段的Tag的vld项作为验证字符串,使用分号分隔验证类型和验证函数
//    格式范例(验证函数参考validator包):
//    `vld:"?;>0&&<10"`
//    `?;>0&&<10`
func AnalyzeField(field *reflect.StructField) (*FieldMetadata, error) {
	var tag = field.Tag.Get("vld")
	tag = strings.TrimSpace(tag)
	if tag == "" && field.Tag != "" {
		var newTag = strings.TrimSpace(string(field.Tag))
		if newTag != "" {
			var firstChar = newTag[0]
			if firstChar == '!' || firstChar == '?' || firstChar == '-' {
				tag = string(field.Tag)
			}
		}
	}
	var fMd = new(FieldMetadata)
	fMd.Name = field.Name
	fMd.Field = field
	switch {
	case strings.HasPrefix(tag, "!"):
		fMd.Kind = FieldKindRequired
	case tag == "" || strings.HasPrefix(tag, "?"):
		fMd.Kind = FieldKindOptional
	case strings.HasPrefix(tag, "-"):
		fMd.Kind = FieldKindIgnore
	default:
		return nil, ErrorInvalidTag.Format(fMd.Name, tag[0]).Error()
	}
	if fMd.Kind == FieldKindOptional || fMd.Kind == FieldKindRequired {
		//获取验证字符串
		var arr = vldReg.FindStringSubmatch(tag)
		if len(arr) == 2 {
			var src = arr[1]
			if len(src) > 0 {
				var vld, err = validator.NewValidator("string", src)
				if err != nil {
					return nil, err
				}
				fMd.Validator = vld
			}
		}
	}
	return fMd, nil
}
