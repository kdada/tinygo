package validator

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
func (this *Parser) Parse() (e error) {
	defer func() {
		//处理表达式解析过程中出现的异常
		var err, ok = recover().(error)
		if ok {
			e = err
		}
	}()
	this.Tree = NewSpaceNode()
	this.expr()
	var t, err = this.Lexer.Prefetch()
	if err != nil {
		return err
	}
	if t.Kind != TokenKindEOF {
		return ErrorInvalidExpr.Format(this.Lexer.Token(t)).Error()
	}
	for this.Tree.Parent() != nil {
		this.Tree = this.Tree.Parent()
	}
	return nil
}

// addSpaceNode 为当前节点添加空间节点
func (this *Parser) addSpaceNode() {
	var node = NewSpaceNode()
	this.Tree.AddChild(node)
	this.Tree = node
}

// addLogicalNode 为当前节点添加逻辑节点
func (this *Parser) addLogicalNode(node SyntaxNode) {
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

// addAndNode 为当前节点添加逻辑与节点
func (this *Parser) addAndNode() {
	this.addLogicalNode(NewAndNode())
}

// addOrNode 为当前节点添加逻辑或节点
func (this *Parser) addOrNode() {
	this.addLogicalNode(NewOrNode())
}

// addFuncNode 为当前节点添加函数节点
func (this *Parser) addFuncNode(name string) {
	var node = NewFuncNode(name)
	this.Tree.AddChild(node)
	this.Tree = node
}

// addParam 给当前函数节点添加参数
func (this *Parser) addParam(t *Token) {
	var node = this.Tree.(*FuncNode)
	node.AddParam(t)
}

// expr 匹配表达式
func (this *Parser) expr() {
	var t, err = this.Lexer.Prefetch()
	if err != nil {
		panic(err)
	}
	if t.Kind == TokenKindEOF {
		return
	}
	if t.Kind == TokenKindLP {
		this.match(TokenKindLP)
		this.addSpaceNode()
		this.expr()
		this.match(TokenKindRP)
		this.extend()
		return
	}
	if t.Kind == TokenKindRelop || t.Kind == TokenKindId || t.Kind == TokenKindRegexp {
		this.function()
		this.extend()
		return
	}
	panic(ErrorInvalidExprHead.Error())
}

// extend 匹配表达式扩展部分
func (this *Parser) extend() {
	var t, err = this.Lexer.Prefetch()
	if err != nil {
		panic(err)
	}
	if t.Kind == TokenKindEOF || t.Kind == TokenKindRP {
		return
	}
	if t.Kind == TokenKindAnd {
		this.match(TokenKindAnd)
		this.addAndNode()
		this.expr()
		return
	}
	if t.Kind == TokenKindOr {
		this.match(TokenKindOr)
		this.addOrNode()
		this.expr()
		return
	}
	panic(ErrorInvalidConnector.Format(this.Lexer.Token(t)).Error())
}

// function 匹配函数
func (this *Parser) function() {
	var t, err = this.Lexer.Prefetch()
	if err != nil {
		panic(err)
	}
	switch t.Kind {
	case TokenKindRelop:
		//匹配关系函数
		this.match(TokenKindRelop)
		this.addFuncNode(t.Value.(string))
		t, err = this.Lexer.Prefetch()
		if err != nil {
			panic(err)
		}
		switch t.Kind {
		case TokenKindLP:
			this.match(TokenKindLP)
			this.optparams()
			this.match(TokenKindRP)
			return
		case TokenKindInteger, TokenKindFloat, TokenKindString:
			this.param()
			return
		}
		panic(ErrorInvalidRelopFuncParams.Format(this.Lexer.Token(t)).Error())
	case TokenKindId:
		this.match(TokenKindId)
		var id = t.StringValue()
		t, err = this.Lexer.Prefetch()
		if err != nil {
			panic(err)
		}
		switch t.Kind {
		case TokenKindLP:
			//匹配函数
			this.addFuncNode(id)
			this.match(TokenKindLP)
			this.optparams()
			this.match(TokenKindRP)
			return
		case TokenKindRelop:
			//匹配命名关系函数
			this.match(TokenKindRelop)
			this.addFuncNode(id + t.StringValue())
			t, err = this.Lexer.Prefetch()
			if err != nil {
				panic(err)
			}
			switch t.Kind {
			case TokenKindLP:
				this.match(TokenKindLP)
				this.optparams()
				this.match(TokenKindRP)
				return
			case TokenKindInteger, TokenKindFloat, TokenKindString:
				this.param()
				return
			}
			panic(ErrorInvalidNamedRelopFuncParams.Format(this.Lexer.Token(t)).Error())
		}
		panic(ErrorInvalidFuncParams.Format(this.Lexer.Token(t)).Error())
	case TokenKindRegexp:
		//匹配正则函数
		this.addFuncNode(t.StringValue())
		this.match(TokenKindRegexp)
		return
	}
	panic(ErrorInvalidFunc.Format(this.Lexer.Token(t)).Error())
}

// optparams 匹配可选的参数列表
func (this *Parser) optparams() {
	var t, err = this.Lexer.Prefetch()
	if err != nil {
		panic(err)
	}
	if t.Kind == TokenKindInteger || t.Kind == TokenKindFloat || t.Kind == TokenKindString {
	FOR:
		for {
			this.param()
			t, err = this.Lexer.Prefetch()
			if err != nil {
				panic(err)
			}
			switch t.Kind {
			case TokenKindSep:
				this.match(TokenKindSep)
			case TokenKindRP, TokenKindEOF:
				break FOR
			default:
				panic(ErrorInvalidParamsList.Format(this.Lexer.Token(t)).Error())
			}
		}
	}
	return
}

// param 匹配单个参数
func (this *Parser) param() {
	var t, err = this.Lexer.Prefetch()
	if err != nil {
		panic(err)
	}
	switch t.Kind {
	case TokenKindInteger:
		this.match(TokenKindInteger)
		this.addParam(t)
		return
	case TokenKindFloat:
		this.match(TokenKindFloat)
		this.addParam(t)
		return
	case TokenKindString:
		this.match(TokenKindString)
		this.addParam(t)
		return
	}
	panic(ErrorInvalidParamType.Format(this.Lexer.Token(t)).Error())
}

// match 匹配指定类型的Token
func (this *Parser) match(kind TokenKind) {
	var t, err = this.Lexer.Fetch()
	if err != nil {
		panic(err)
	}
	if t.Kind == kind {
		return
	}
	panic(ErrorUnmatchedToken.Format(kind, t.Kind).Error())
}
