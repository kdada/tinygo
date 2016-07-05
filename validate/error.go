package validate

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
	ErrorInvalidValidatorCreator Error = "ErrorInvalidValidatorCreator(V10000):无效的Validator创建器"
	ErrorInvalidValidatorKind    Error = "ErrorInvalidValidatorKind(V10010):无效的Validator类型(%s)"
	ErrorUnmatchRegBoundary      Error = "ErrorUnmatchRegBoundary(V10020):无效的正则表达式边界,缺少/"
	ErrorUnmatchStringBoundary   Error = "ErrorUnmatchStringBoundary(V10021):无效的字符串边界,缺少'"
	ErrorUnmatchAnd              Error = "ErrorUnmatchAnd(V10030):逻辑与必须有两个连续的&符号"
	ErrorUnmatchOr               Error = "ErrorUnmatchOr(V10040):逻辑或必须有两个连续的|符号"
	ErrorInvalidNumber           Error = "ErrorInvalidNumber(V10050):负号(-)后必须存在数字"
)
