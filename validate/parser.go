package validate

import "errors"

type ScopeKind string

const (
	ScopeKindSpace ScopeKind = "空范围"   //()范围
	ScopeKindAnd   ScopeKind = "And范围" //&&范围
	ScopeKindOr    ScopeKind = "Or范围"  //||范围
	ScopeKindFunc  ScopeKind = "函数范围"  //方法
)

// 范围节点
type ScopeNode struct {
	Kind   ScopeKind
	Parent *ScopeNode
	Left   *ScopeNode
	Right  *ScopeNode
	Value  interface{}
}

func (this *ScopeNode) AddLeft(node *ScopeNode) {
	node.Parent = this
	this.Left = node
}

func (this *ScopeNode) AddRight(node *ScopeNode) {
	node.Parent = this
	this.Right = node
}

func (this *ScopeNode) HasLeft() bool {
	return this.Left != nil
}
func (this *ScopeNode) HasRight() bool {
	return this.Right != nil
}

// NewScopeNode 创建范围节点
func NewScopeNode(kind ScopeKind) *ScopeNode {
	return &ScopeNode{kind, nil, nil, nil, nil}
}

// Parse 生成语法树
func Parse(tokenizor *Tokenizor) (*ScopeNode, error) {
	var root = NewScopeNode(ScopeKindSpace)
	var cur = root
	for {
		var token, err = tokenizor.Fetch()
		if err == EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		switch token.Kind {
		case TokenKindLeft:
			{
				var n = NewScopeNode(ScopeKindSpace)
				if !cur.HasLeft() {
					cur.AddLeft(n)
				} else {
					cur.AddRight(n)
				}
				cur = n
			}
		case TokenKindRight:
			{
				cur = cur.Parent
				if cur == nil {
					return nil, errors.New("括号数量不匹配")
				}
			}
		case TokenKindAnd:
			{
				if !cur.HasLeft() {
					return nil, errors.New("&&前必须有表达式")
				}
				if cur.HasRight() {
					if cur.Parent == nil {
						return nil, errors.New("括号数量不匹配")
					}
					var node = NewScopeNode(ScopeKindAnd)
					var p = cur.Parent
					node.AddLeft(cur)
					if p.Left == cur {
						p.AddLeft(node)
					} else {
						p.AddRight(node)
					}
					cur = node
				} else {
					cur.Kind = ScopeKindAnd
				}
			}
		case TokenKindOr:
			{
				if !cur.HasLeft() {
					return nil, errors.New("||前必须有表达式")
				}
				if cur.HasRight() {
					if cur.Parent == nil {
						return nil, errors.New("括号数量不匹配")
					}
					var node = NewScopeNode(ScopeKindOr)
					var p = cur.Parent
					node.AddLeft(cur)
					if p.Left == cur {
						p.AddLeft(node)
					} else {
						p.AddRight(node)
					}
					cur = node
				} else {
					cur.Kind = ScopeKindOr
				}
			}
		case TokenKindReg:
			{
				var n = NewScopeNode(ScopeKindFunc)
				n.Value = token.Value
				if !cur.HasLeft() {
					cur.AddLeft(n)
				} else if !cur.HasRight() {
					cur.AddRight(n)
				}
			}
		case TokenKindFunc:
			{
				var n = NewScopeNode(ScopeKindFunc)
				var str = token.Value
				var p, err = tokenizor.Prefetch()
				if err == nil && p.Kind != TokenKindAnd && p.Kind != TokenKindOr && p.Kind != TokenKindRight {
					tokenizor.Fetch()
					if p.Kind != TokenKindLeft {
						str += " " + p.Value
					} else {
						var lastIsParam = false
						var rightParenthesis = false
						for {
							p, err = tokenizor.Fetch()
							if err != nil {
								break
							}
							if p.Kind == TokenKindRight {
								rightParenthesis = true
								break
							}
							if p.Kind != TokenKindInt || p.Kind != TokenKindFloat || p.Kind != TokenKindString || p.Kind != TokenKindSeparator {
								break
							}
							if lastIsParam {
								lastIsParam = false
								continue
							}
							lastIsParam = true
							str += " " + p.Value
						}
						if !rightParenthesis {
							err = errors.New("函数右括号缺失")
						}
					}
				}
				if err != nil {
					return nil, err
				}
				n.Value = str
				if !cur.HasLeft() {
					cur.AddLeft(n)
				} else if !cur.HasRight() {
					cur.AddRight(n)
				}
			}
		}
	}
	if cur != root {
		return nil, errors.New("括号数量不匹配")
	}
	return root, nil
}
