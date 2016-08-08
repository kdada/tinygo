package meta

import (
	"errors"
	"fmt"
)

// 错误信息
type Error string

// Format 格式化错误信息并生成新的错误信息
func (this Error) Format(data ...interface{}) Error {
	return Error(fmt.Sprintf(string(this), data...))
}

// Error 生成error类型
func (this Error) Error() error {
	return errors.New(string(this))
}

// String 返回错误字符串描述
func (this Error) String() string {
	return string(this)
}

// 错误码
const (
	ErrorParamNotFunc   Error = "ErrorParamNotFunc(M10010):%s不是函数"
	ErrorParamNotStruct Error = "ErrorParamNotStruct(M10020):%s不是结构体类型"
	ErrorInvalidTag     Error = "ErrorInvalidTag(M10030):字段(%s)的vld验证字符串必须为空或者第一个字符为!(必须),?(可选),-(忽略),首个字符不能是(%s)"
	ErrorRequiredField  Error = "ErrorRequiredField(M10040):字段(%s)的值不存在"
	ErrorFieldNotValid  Error = "ErrorFieldNotValid(M10041):字段(%s)无法通过校验"
	ErrorFieldsNotValid Error = "ErrorFieldsNotValid(M10042):字段(%s)的第(%d)个值无法通过校验"

	ErrorMustNotBeInterface Error = "ErrorMustNotBeInterface(M10050):指定参数不能是接口类型"
	ErrorNotAssignable      Error = "ErrorNotAssignable(M10051):指定类型(%s)不能赋值给类型(%s)"
)
