package validator

import "errors"

/*
词法:
any			-> .*
letter		-> [A-Za-z]
digit		-> [0-9]
number		-> digit+
integer		-> ('+'|'-'|ε)number
float		-> integer'.'number
string		-> '\''any'\''
relop		-> '>' | '>=' | '<' | '<=' | '==' | '!='
id			-> letter(letter|digit)*
regexp		-> '/'any'/'
and			-> '&&'
or			-> '||'
lp			-> '('
rp			-> ')'
sep			-> ','
*/

// 标记类型
type TokenKind string

const (
	TokenKindEOF     TokenKind = "TokenKindEOF"     //结束标记
	TokenKindInteger TokenKind = "TokenKindInteger" //整数
	TokenKindFloat   TokenKind = "TokenKindFloat"   //浮点数
	TokenKindString  TokenKind = "TokenKindString"  //字符串
	TokenKindRelop   TokenKind = "TokenKindRelop"   //关系运算符
	TokenKindId      TokenKind = "TokenKindId"      //id
	TokenKindRegexp  TokenKind = "TokenKindRegexp"  //正则表达式
	TokenKindAnd     TokenKind = "TokenKindAnd"     //逻辑与
	TokenKindOr      TokenKind = "TokenKindOr"      //逻辑或
	TokenKindLP      TokenKind = "TokenKindLP"      //左括号
	TokenKindRP      TokenKind = "TokenKindRP"      //右括号
	TokenKindSep     TokenKind = "TokenKindSep"     //逗号分隔符
)

// 标记
type Token struct {
	Kind TokenKind //标记类型
	Pos  int       //起始位置
	Len  int       //长度
}

// NewToken 创建标记
func NewToken(kind TokenKind, pos int, length int) *Token {
	return &Token{kind, pos, length}
}

// 词法分析器
type Lexer struct {
	Data []rune //字符串
	Pos  int    //下一个标记开始的位置
}

// NewLexer 创建词法解析器
func NewLexer(src string) *Lexer {
	return &Lexer{[]rune(src), 0}
}
func (this *Lexer) Token(t *Token) string {
	if t.Len <= 0 {
		return ""
	}
	return string(this.Data[t.Pos : t.Pos+t.Len])
}

// Fetch 获取下一个标记,已经没有有效标记的时候会返回TokenKindEOF类型的标记
func (this *Lexer) Fetch() (*Token, error) {
	var t, e = this.Prefetch()
	if e == nil {
		this.Pos += t.Len
	}
	return t, e
}

// Prefetch 预获取下一个标记,已经没有有效标记的时候会返回TokenKindEOF类型的标记
func (this *Lexer) Prefetch() (*Token, error) {
	for i := this.Pos; i < len(this.Data); i++ {
		var c = this.Data[i]
		if c != ' ' {
			this.Pos = i
			switch c {
			case '(':
				return NewToken(TokenKindLP, this.Pos, 1), nil
			case ')':
				return NewToken(TokenKindRP, this.Pos, 1), nil
			case ',':
				return NewToken(TokenKindSep, this.Pos, 1), nil
			case '&':
				return this.and()
			case '|':
				return this.or()
			case '/':
				return this.regexp()
			case '>', '<', '!', '=':
				return this.relop()
			case '+', '-', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
				return this.number()
			case '\'':
				return this.string()
			default:
				if (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') {
					return this.id()
				}
			}
		}
	}
	if this.Pos >= len(this.Data) {
		return NewToken(TokenKindEOF, this.Pos, 0), nil
	}
	return nil, errors.New("无效的字符")
}

// and 用于识别&&
func (this *Lexer) and() (*Token, error) {
	if this.Pos+1 < len(this.Data) && this.Data[this.Pos+1] == '&' {
		return NewToken(TokenKindAnd, this.Pos, 2), nil
	}
	return nil, errors.New("逻辑与缺少&字符")
}

// or 用于识别||
func (this *Lexer) or() (*Token, error) {
	if this.Pos+1 < len(this.Data) && this.Data[this.Pos+1] == '|' {
		return NewToken(TokenKindOr, this.Pos, 2), nil
	}
	return nil, errors.New("逻辑或缺少|字符")
}

// regexp 用于识别正则表达式
func (this *Lexer) regexp() (*Token, error) {
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
				return NewToken(TokenKindRegexp, this.Pos, length), nil
			}
		}
	}
	return nil, errors.New("正则表达式缺少结束标记")
}

// relop 用于识别关系表达式
func (this *Lexer) relop() (*Token, error) {
	var fc = this.Data[this.Pos]
	var sc = rune(0)
	if this.Pos+1 < len(this.Data) {
		sc = this.Data[this.Pos+1]
	}
	var length = 1
	if sc == '=' {
		length = 2
	}
	if length == 1 && (fc == '!' || fc == '=') {
		return nil, errors.New("关系表达式不完整")
	}
	return NewToken(TokenKindRelop, this.Pos, length), nil
}

// number 用于识别整数和浮点数
func (this *Lexer) number() (*Token, error) {
	var c = this.Data[this.Pos]
	var startWithNum = true
	if !(c >= '0' && c <= '9') {
		startWithNum = false
	}
	var i = this.Pos + 1
	var dotPos = -1
	for ; i < len(this.Data); i++ {
		var b = this.Data[i]
		if b == '.' {
			if i == this.Pos+1 && !startWithNum {
				return nil, errors.New("小数点不能紧跟+/-")
			}
			if dotPos > 0 {
				return nil, errors.New("无效的浮点数,存在多个小数点")
			}
			dotPos = i
			continue
		}
		if !(b >= '0' && b <= '9') {
			break
		}
	}
	var length = i - this.Pos
	if length <= 1 && !startWithNum {
		return nil, errors.New("数值常量中必须包含数字")
	}
	if dotPos >= 0 {
		return NewToken(TokenKindFloat, this.Pos, length), nil
	}
	return NewToken(TokenKindInteger, this.Pos, length), nil
}

// string 用于识别字符串
func (this *Lexer) string() (*Token, error) {
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
				return NewToken(TokenKindString, this.Pos, i-this.Pos+1), nil
			}
		}
	}
	return nil, errors.New("字符串缺少结束标记")
}

// id 用于识别标识符
func (this *Lexer) id() (*Token, error) {
	var i = this.Pos + 1
	for ; i < len(this.Data); i++ {
		var c = this.Data[i]
		if !((c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9')) {
			break
		}
	}
	return NewToken(TokenKindId, this.Pos, i-this.Pos), nil
}
