package validator

import (
	"errors"
	"fmt"
)

/*
语法:
param		-> integer
			|  float
			|  string
optparams	-> param(sep param)*
			|  ε
function	-> relop param
			|  relop lp optparams rp
			|  id relop param
			|  id relop lp optparams rp
			|  id lp optparams rp
			|  regexp
extend		-> and expr
			|  or expr
			|  ε
expr		-> lp expr rp extend
			|  function extend


FIRST(param)	= {integer,float,string}
FIRST(optparams) = {integer,float,string,ε}
FIRST(function) = {relop,id,regexp}
FIRST(extend) = {and,or,ε}
FIRST(expr) = {lp,relop,id,regexp}
*/

// 语法分析器
type Parser struct {
	Lexer *Lexer
	Tree  SyntaxNode
}

// NewParser 创建语法分析器
func NewParser(l *Lexer) *Parser {
	return &Parser{l, nil}
}

// Parse 解析成语法树
func (this *Parser) Parse() error {
	this.Tree = NewSpaceNode()
	var err = this.expr()
	if err != nil {
		return err
	}
	var t, err2 = this.Lexer.Prefetch()
	if err2 != nil {
		return err2
	}
	if t.Kind != TokenKindEOF {
		fmt.Println("表达式尾有无效的内容")
		return errors.New("表达式尾有无效的内容")
	}
	for this.Tree.Parent() != nil {
		this.Tree = this.Tree.Parent()
	}
	return nil
}

func (this *Parser) addSpaceNode() {
	var node = NewSpaceNode()
	this.Tree.AddChild(node)
	this.Tree = node
}

func (this *Parser) addRelopNode(node SyntaxNode) {
	var p = this.Tree.Parent()
	if p.Left() == this.Tree {
		p.ChangeKind(node.Kind())
		this.Tree = p
	} else {
		var gp = p.Parent()
		if gp == nil {
			gp = NewSpaceNode()
			gp.AddChild(p)
		}
		if gp.Left() == p {
			gp.ChangeKind(node.Kind())
			this.Tree = gp
		} else {
			node.AddChild(p)
			gp.SetRight(node)
			this.Tree = node
		}
	}

}

func (this *Parser) addAndNode() {
	this.addRelopNode(NewAndNode())
}

func (this *Parser) addOrNode() {
	this.addRelopNode(NewOrNode())
}

func (this *Parser) addFuncNode(name string) {
	var node = NewFuncNode(name)
	this.Tree.AddChild(node)
	this.Tree = node
}

func (this *Parser) addParam(t *Token) {
	var node = this.Tree.(*FuncNode)
	node.AddParam(t)
}

func (this *Parser) expr() error {
	var t, err = this.Lexer.Prefetch()
	if err != nil {
		return err
	}
	if t.Kind == TokenKindEOF {
		return nil
	}
	if t.Kind == TokenKindLP {
		this.match(TokenKindLP)
		this.addSpaceNode()
		this.expr()
		this.match(TokenKindRP)
		this.extend()
		return nil
	}
	if t.Kind == TokenKindRelop || t.Kind == TokenKindId || t.Kind == TokenKindRegexp {
		this.function()
		this.extend()
		return nil
	}
	fmt.Println("表达式语法错误,无法识别的表达式")
	return errors.New("表达式语法错误,无法识别的表达式")
}

func (this *Parser) extend() error {
	var t, err = this.Lexer.Prefetch()
	if err != nil {
		return err
	}
	if t.Kind == TokenKindEOF || t.Kind == TokenKindRP {
		return nil
	}
	if t.Kind == TokenKindAnd {
		this.match(TokenKindAnd)
		this.addAndNode()
		this.expr()
		return nil
	}
	if t.Kind == TokenKindOr {
		this.match(TokenKindOr)
		this.addOrNode()
		this.expr()
		return nil
	}
	fmt.Println("扩展语法错误,必须使用&&或||连接表达式")
	return errors.New("扩展语法错误,必须使用&&或||连接表达式")
}

func (this *Parser) function() error {
	var t, err = this.Lexer.Prefetch()
	if err != nil {
		return err
	}
	switch t.Kind {
	case TokenKindRelop:
		this.match(TokenKindRelop)
		this.addFuncNode(this.Lexer.Token(t))
		t, err = this.Lexer.Prefetch()
		if err != nil {
			return err
		}
		switch t.Kind {
		case TokenKindLP:
			this.match(TokenKindLP)
			this.optparams()
			this.match(TokenKindRP)
			return nil
		case TokenKindInteger, TokenKindFloat, TokenKindString:
			this.param()
			return nil
		}
		fmt.Println("关系运算符函数参数错误")
		return errors.New("关系运算符函数参数错误")
	case TokenKindId:
		this.match(TokenKindId)
		var id = this.Lexer.Token(t)
		t, err = this.Lexer.Prefetch()
		if err != nil {
			return err
		}
		switch t.Kind {
		case TokenKindLP:
			this.addFuncNode(id)
			this.match(TokenKindLP)
			this.optparams()
			this.match(TokenKindRP)
			return nil
		case TokenKindRelop:
			this.match(TokenKindRelop)
			this.addFuncNode(id + this.Lexer.Token(t))
			t, err = this.Lexer.Prefetch()
			if err != nil {
				return err
			}
			switch t.Kind {
			case TokenKindLP:
				this.match(TokenKindLP)
				this.optparams()
				this.match(TokenKindRP)
				return nil
			case TokenKindInteger, TokenKindFloat, TokenKindString:
				this.param()
				return nil
			}
			fmt.Println("命名关系函数参数错误")
			return errors.New("命名关系函数参数错误")
		}
		fmt.Println("函数参数错误")
		return errors.New("函数参数错误")
	case TokenKindRegexp:
		this.addFuncNode(this.Lexer.Token(t))
		this.match(TokenKindRegexp)
		return nil
	}
	fmt.Println("无法识别的函数")
	return errors.New("无法识别的函数")
}

func (this *Parser) optparams() error {
	var t, err = this.Lexer.Prefetch()
	if err != nil {
		return err
	}
	if t.Kind == TokenKindInteger || t.Kind == TokenKindFloat || t.Kind == TokenKindString {
	FOR:
		for {
			this.param()
			t, err = this.Lexer.Prefetch()
			if err != nil {
				return err
			}
			switch t.Kind {
			case TokenKindSep:
				this.match(TokenKindSep)
			case TokenKindRP, TokenKindEOF:
				break FOR
			default:
				fmt.Println("参数列表错误,使用了不正确的分隔符或缺少右括号")
				return errors.New("参数列表错误,使用了不正确的分隔符")
			}
		}
	}
	return nil
}

func (this *Parser) param() error {
	var t, err = this.Lexer.Prefetch()
	if err != nil {
		return err
	}
	switch t.Kind {
	case TokenKindInteger:
		this.match(TokenKindInteger)
		this.addParam(t)
		return nil
	case TokenKindFloat:
		this.match(TokenKindFloat)
		this.addParam(t)
		return nil
	case TokenKindString:
		this.match(TokenKindString)
		this.addParam(t)
		return nil
	}
	fmt.Println("参数类型错误")
	return errors.New("参数类型错误")
}

func (this *Parser) match(kind TokenKind) error {
	var t, err = this.Lexer.Fetch()
	if err != nil {
		return err
	}
	if t.Kind == kind {
		return nil
	}
	fmt.Println("无法匹配的TokenKind")
	return errors.New("无法匹配的TokenKind")
}
