package validate

import (
	"errors"
	"strings"
)

// 结束标记错误
var EOF = errors.New("EOF")

// Token类型
type TokenKind byte

const (
	TokenKindLeft      TokenKind = iota //(
	TokenKindRight                      //)
	TokenKindAnd                        //&&
	TokenKindOr                         //||
	TokenKindSeparator                  //,
	TokenKindFunc                       //函数标识符
	TokenKindInt                        //数字标识符
	TokenKindFloat                      //浮点数标识符
	TokenKindString                     //字符串标识符
	TokenKindReg                        //正则表达式标识符
)

// Token
type Token struct {
	Kind  TokenKind //Token类型
	Pos   int       //起始位置
	Len   int       //长度
	Value string    //值
}

// NewToken 创建新的Token
func NewToken(kind TokenKind, pos int, length int, value string) *Token {
	return &Token{kind, pos, length, value}
}

// Token化工具
type Tokenizor struct {
	Data []byte
	Pos  int
}

// NewTokenizor 创建Token工具
func NewTokenizor(src []byte) *Tokenizor {
	return &Tokenizor{src, 0}
}

// Fetch 获取下一个Token
func (this *Tokenizor) Fetch() (*Token, error) {
	var t, e = this.Prefetch()
	if e == nil {
		this.Pos += t.Len
	}
	return t, e
}

// Prefetch 预获取下一个Token,不修改读取位置
func (this *Tokenizor) Prefetch() (*Token, error) {
	for i := this.Pos; i < len(this.Data); i++ {
		if this.Data[i] != ' ' {
			this.Pos = i
			switch this.Data[i] {
			case '(':
				return NewToken(TokenKindLeft, this.Pos, 1, "("), nil
			case ')':
				return NewToken(TokenKindRight, this.Pos, 1, ")"), nil
			case ',':
				return NewToken(TokenKindSeparator, this.Pos, 1, ","), nil
			case '/':
				return this.reg()
			case '&':
				return this.and()
			case '|':
				return this.or()
			default:
				return this.id()
			}
		}
	}
	return nil, EOF
}

// reg 读取一个正则
func (this *Tokenizor) reg() (*Token, error) {
	var lastCharIsEscape = false
	for i := this.Pos + 1; i < len(this.Data); i++ {
		if lastCharIsEscape {
			lastCharIsEscape = false
			continue
		}
		var b = this.Data[i]
		switch b {
		case '\\':
			{
				lastCharIsEscape = true
			}
		case '/':
			{
				var length = i - this.Pos + 1
				return NewToken(TokenKindReg, this.Pos, length, strings.Replace(string(this.Data[this.Pos+1:this.Pos+length-1]), "\\/", "/", -1)), nil
			}
		}
	}
	return nil, ErrorUnmatchRegBoundary.Error()
}

// reg 读取一个逻辑与
func (this *Tokenizor) and() (*Token, error) {
	if this.Pos+1 < len(this.Data) && this.Data[this.Pos+1] == '&' {
		return NewToken(TokenKindAnd, this.Pos, 2, "&&"), nil
	}
	return nil, ErrorUnmatchAnd.Error()
}

// or 读取一个逻辑或
func (this *Tokenizor) or() (*Token, error) {
	if this.Pos+1 < len(this.Data) && this.Data[this.Pos+1] == '|' {
		return NewToken(TokenKindOr, this.Pos, 2, "||"), nil
	}
	return nil, ErrorUnmatchOr.Error()
}

// id 读取一个标识符
func (this *Tokenizor) id() (*Token, error) {
	var b = this.Data[this.Pos]
	if b == '\'' {
		return this.string()
	}
	if b == '-' {
		if this.Pos+1 >= len(this.Data) {
			return nil, ErrorInvalidNumber.Error()
		}
		var n = this.Data[this.Pos+1]
		if !this.isNumber(n) {
			return nil, ErrorInvalidNumber.Error()
		}
		return this.int()
	}
	if this.isNumber(b) {
		return this.int()
	}
	return this.function()
}

// int 读取一个整数
func (this *Tokenizor) int() (*Token, error) {
	var i = this.Pos + 1
	for ; i < len(this.Data); i++ {
		var b = this.Data[i]
		if b == '.' {
			return this.float()
		}
		if !this.isNumber(b) {
			break
		}
	}
	return NewToken(TokenKindInt, this.Pos, i-this.Pos, string(this.Data[this.Pos:i])), nil
}

// float 读取一个浮点数
func (this *Tokenizor) float() (*Token, error) {
	var i = this.Pos + 1
	for ; i < len(this.Data); i++ {
		var b = this.Data[i]
		if !(this.isNumber(b) || b == '.') {
			break
		}
	}
	return NewToken(TokenKindFloat, this.Pos, i-this.Pos, string(this.Data[this.Pos:i])), nil
}

// string 读取一个用单引号包裹的字符串
func (this *Tokenizor) string() (*Token, error) {
	var lastCharIsEscape = false
	for i := this.Pos + 1; i < len(this.Data); i++ {
		if lastCharIsEscape {
			lastCharIsEscape = false
			continue
		}
		var b = this.Data[i]
		switch b {
		case '\\':
			{
				lastCharIsEscape = true
			}
		case '\'':
			{
				return NewToken(TokenKindString, this.Pos, i-this.Pos+1, strings.Replace(string(this.Data[this.Pos+1:i]), "\\'", "'", -1)), nil
			}
		}
	}
	return nil, ErrorUnmatchStringBoundary.Error()

}

// function 读取一个函数名,函数名中不能包含保留字符
func (this *Tokenizor) function() (*Token, error) {
	var i = this.Pos + 1
	var last = this.Data[this.Pos]
	for ; i < len(this.Data); i++ {
		var b = this.Data[i]
		//函数名不能包含(和-以及'
		//上一个字符不是字母和数字时,当前字符不能为数字
		if this.isReserved(b) || (!(this.isLetter(last) || this.isNumber(last)) && this.isNumber(b)) {
			break
		}
		last = b
	}
	return NewToken(TokenKindFunc, this.Pos, i-this.Pos, strings.Replace(string(this.Data[this.Pos:i]), " ", "", -1)), nil
}

// isLetter 是否是字母
func (this *Tokenizor) isLetter(b byte) bool {
	return (b >= 'a' && b <= 'z') || (b >= 'A' && b <= 'Z')
}

// isNumber 是否是数字
func (this *Tokenizor) isNumber(b byte) bool {
	return b >= '0' && b <= '9'
}

// isReserved 是否是保留字符
func (this *Tokenizor) isReserved(b byte) bool {
	var str = "()&|\\'/,.-"
	return strings.ContainsRune(str, rune(b))
}
