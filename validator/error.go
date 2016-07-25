package validator

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
	ErrorInvalidChar         Error = "ErrorInvalidChar(V11000):无效的字符(位置:%d,%s)"
	ErrorInvalidLogicalAnd   Error = "ErrorInvalidLogicalAnd(V11010):逻辑与的形式必须是&&(位置:%d)"
	ErrorInvalidLogicalOr    Error = "ErrorInvalidLogicalOr(V11020):逻辑或的形式必须是||(位置:%d)"
	ErrorInvalidRelop        Error = "ErrorInvalidRelop(V11030):关系运算符错误,缺少符号=(位置:%d)"
	ErrorUnmatchEnding       Error = "ErrorUnmatchEnding(V11040):字符序列缺少结束标记(位置:%d,%s)"
	ErrorInvalidNumberFormat Error = "ErrorInvalidNumberFormat(V11050):数值格式错误,+/-之后必须是数字(位置:%d)"
	ErrorInvalidFloat        Error = "ErrorInvalidFloat(V11051):浮点格式错误,浮点数中不能包含两个小数点(位置:%d)"

	ErrorInvalidExpr                 Error = "ErrorInvalidExpr(V12000):无效的表达式,表达式尾部含有未被解析的内容(位置:%d,%s)"
	ErrorInvalidExprHead             Error = "ErrorInvalidExprHead(V12001):无效的表达式,表达式的第一个标记必须是左括号或函数(位置:%d)"
	ErrorInvalidConnector            Error = "ErrorInvalidConnector(V12010):必须使用&&或||连接表达式或函数,而不是(位置:%d,%s)"
	ErrorInvalidRelopFuncParams      Error = "ErrorInvalidRelopFuncParams(V12020):关系函数参数列表错误(位置:%d,%s)"
	ErrorInvalidNamedRelopFuncParams Error = "ErrorInvalidNamedRelopFuncParams(V12021):命名关系函数参数列表错误(位置:%d,%s)"
	ErrorInvalidFuncParams           Error = "ErrorInvalidFuncParams(V12022):函数参数列表错误(位置:%d,%s)"
	ErrorInvalidFunc                 Error = "ErrorInvalidFunc(V12030):无法识别的函数(位置:%d,%s)"
	ErrorInvalidParamsList           Error = "ErrorInvalidParamsList(V12040):参数列表中存在无效的分隔符(位置:%d,%s)"
	ErrorInvalidParamType            Error = "ErrorInvalidParamType(V12050):参数列表中使用了错误的参数(位置:%d,%s)"
	ErrorUnmatchedToken              Error = "ErrorUnmatchedToken(V12060):类型匹配错误,需要的类型为(%s),实际类型为(位置:%d,%s)"
)
