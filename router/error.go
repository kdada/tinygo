package router

import (
	"errors"
	"fmt"
)

// 错误信息
type RouterError string

// 错误码
const (
	RouterErrorRegexpParseError    RouterError = "T10010:RouterErrorRegexpParseError,无效的正则表达式(%s)"
	RouterErrorRegexpNoneError     RouterError = "T10020:RouterErrorRegexpNoneError,字符串(%s)不包含正则表达式"
	RouterErrorRegexpFormatError   RouterError = "T10030:RouterErrorRegexpFormatError,字符串(%s)格式错误"
	RouterErrorRegexpNotMatchError RouterError = "T10040:RouterErrorRegexpNotMatchError,正则表达式匹配失败"
)

// Format 格式化错误信息并生成新的错误信息
func (this RouterError) Format(data ...interface{}) RouterError {
	return RouterError(fmt.Sprintf(string(this), data...))
}

// Error 生成error类型
func (this RouterError) Error() error {
	return errors.New(string(this))
}
