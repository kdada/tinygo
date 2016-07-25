package validator

import "strconv"

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
	Kind  TokenKind   //标记类型
	Pos   int         //起始位置
	Len   int         //长度
	Value interface{} //Token值,TokenKindInteger->int64,TokenKindFloat->float64,其他->string
}

// StringValue 返回字符串类型的值,使用前确保Value为string类型
func (this *Token) StringValue() string {
	return this.Value.(string)
}

// Int64Value 返回整数类型的值,使用前确保Value为int64类型
func (this *Token) Int64Value() int64 {
	return this.Value.(int64)
}

// Float64Value 返回浮点类型的值,使用前确保Value为float64类型
func (this *Token) Float64Value() float64 {
	return this.Value.(float64)
}

// NewToken 创建标记
func NewToken(kind TokenKind, pos int, length int, value interface{}) *Token {
	return &Token{kind, pos, length, value}
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
				return NewToken(TokenKindLP, this.Pos, 1, "("), nil
			case ')':
				return NewToken(TokenKindRP, this.Pos, 1, ")"), nil
			case ',':
				return NewToken(TokenKindSep, this.Pos, 1, ","), nil
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
				return nil, ErrorInvalidChar.Format(string(c)).Error()
			}
		}
	}
	return NewToken(TokenKindEOF, this.Pos, 0, nil), nil
}

// Token 获取Token字符串
func (this *Lexer) Token(t *Token) string {
	if t.Len <= 0 {
		return ""
	}
	return string(this.Data[t.Pos : t.Pos+t.Len])
}

// and 用于识别&&
func (this *Lexer) and() (*Token, error) {
	if this.Pos+1 < len(this.Data) && this.Data[this.Pos+1] == '&' {
		return NewToken(TokenKindAnd, this.Pos, 2, "&&"), nil
	}
	return nil, ErrorInvalidLogicalAnd.Error()
}

// or 用于识别||
func (this *Lexer) or() (*Token, error) {
	if this.Pos+1 < len(this.Data) && this.Data[this.Pos+1] == '|' {
		return NewToken(TokenKindOr, this.Pos, 2, "||"), nil
	}
	return nil, ErrorInvalidLogicalOr.Error()
}

// regexp 用于识别正则表达式
func (this *Lexer) regexp() (*Token, error) {
	var str, pos, err = this.unescaped(this.Pos+1, '/')
	if err != nil {
		return nil, err
	}
	return NewToken(TokenKindRegexp, this.Pos, pos-this.Pos+1, str), nil
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
		return nil, ErrorInvalidRelop.Error()
	}
	return NewToken(TokenKindRelop, this.Pos, length, string(this.Data[this.Pos:this.Pos+length])), nil
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
				return nil, ErrorInvalidNumberFormat.Error()
			}
			if dotPos > 0 {
				return nil, ErrorInvalidFloat.Error()
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
		return nil, ErrorInvalidNumberFormat.Error()
	}
	var s = string(this.Data[this.Pos : this.Pos+length])
	if dotPos >= 0 {
		var f, err = strconv.ParseFloat(s, 64)
		if err != nil {
			return nil, err
		}
		return NewToken(TokenKindFloat, this.Pos, length, f), nil
	}
	var num, err = strconv.ParseInt(s, 10, 64)
	if err != nil {
		return nil, err
	}
	return NewToken(TokenKindInteger, this.Pos, length, num), nil
}

// string 用于识别字符串
func (this *Lexer) string() (*Token, error) {
	var str, pos, err = this.unescaped(this.Pos+1, '\'')
	if err != nil {
		return nil, err
	}
	return NewToken(TokenKindString, this.Pos, pos-this.Pos+1, str), nil
}

// unescaped 用于检索字符串直到遇到end字符,通过双写end字符可对end字符转义.
//  return:(解析后的字符串,出现end字符的位置,解析过程中出现的错误)
func (this *Lexer) unescaped(startPos int, end rune) (string, int, error) {
	var str = make([]rune, 0, 10)
	for i := startPos; i < len(this.Data); i++ {
		var b = this.Data[i]
		if b == end {
			if !(len(this.Data) > i+1 && this.Data[i+1] == end) {
				return string(str), i, nil
			}
			i++
		}
		str = append(str, b)
	}
	return "", startPos, ErrorUnmatchEnding.Format(end).Error()
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
	return NewToken(TokenKindId, this.Pos, i-this.Pos, string(this.Data[this.Pos:i])), nil
}
