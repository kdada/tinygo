package validator

import (
	"reflect"
	"regexp"
)

// 验证器函数表
var funcs = make(map[string]ValidatorFunc)

// 函数与参数列表的分隔符
const sep = ":"

var funcNameReg = regexp.MustCompile("^([a-zA-Z][a-zA-Z0-9]*)?(>|>=|<|<=|==|!=)?$")

// RegisterFunc 注册验证函数,验证函数的第一个参数必须是string,并且其他参数只能是int64,float64,string三种类型,返回值必须是bool类型
func RegisterFunc(name string, f interface{}) error {
	if !funcNameReg.MatchString(name) {
		return ErrorInvalidFuncName.Format(name).Error()
	}
	var v = reflect.ValueOf(f)
	if v.Kind() != reflect.Func {
		return ErrorMustBeFunc.Format(v.Kind().String()).Error()
	}
	var vType = v.Type()
	if vType.NumIn() <= 0 {
		return ErrorIncorrectParamList.Format(name).Error()
	}
	var firstParam = vType.In(0)
	if firstParam.Kind() != reflect.String {
		return ErrorFirstParamMustBeString.Format(name).Error()
	}
	var newName = name + sep
	for i := 1; i < vType.NumIn(); i++ {
		var t = vType.In(i)
		var n, ok = CheckType(t.Kind())
		if !ok {
			return ErrorIncorrectParamType.Format(i, t.Kind().String()).Error()
		}
		newName += n
	}
	funcs[newName] = &NamedFunc{name, &v}
	return nil
}

// CheckType 检查类型是否符合参数要求
func CheckType(k reflect.Kind) (string, bool) {
	switch k {
	case reflect.Int64:
		return "I", true
	case reflect.Float64:
		return "F", true
	case reflect.String:
		return "S", true
	}
	return "", false
}

// 验证器函数
type ValidatorFunc interface {
	Validate(str string, optparams ...interface{}) bool
}

// 命名函数
type NamedFunc struct {
	Name string         //函数名称
	Func *reflect.Value //函数
}

// Validate 根据传入参数进行验证,并返回验证结果
func (this *NamedFunc) Validate(str string, optparams ...interface{}) bool {
	var params = make([]reflect.Value, len(optparams)+1)
	params[0] = reflect.ValueOf(str)
	for i, v := range optparams {
		params[i+1] = reflect.ValueOf(v)
	}
	var results = this.Func.Call(params)
	return results[0].Bool()
}

// 正则函数
type RegFunc struct {
	Exp *regexp.Regexp
}

// NewRegFunc 创建正则函数
func NewRegFunc(reg string) (ValidatorFunc, error) {
	var r, err = regexp.Compile(reg)
	if err != nil {
		return nil, err
	}
	return &RegFunc{r}, nil
}

// Validate 根据传入参数进行验证,并返回验证结果
func (this *RegFunc) Validate(str string, optparams ...interface{}) bool {
	return this.Exp.MatchString(str)
}
