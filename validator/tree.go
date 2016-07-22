package validator

type NodeKind string

const (
	NodeKindSpace NodeKind = "NodeKindSpace"
	NodeKindAnd   NodeKind = "NodeKindAnd"
	NodeKindOr    NodeKind = "NodeKindOr"
	NodeKindFunc  NodeKind = "NodeKindFunc"
)

// 语法节点接口
type SyntaxNode interface {
	Kind() NodeKind
	ChangeKind(kind NodeKind)
	Parent() SyntaxNode
	SetParent(n SyntaxNode)
	Left() SyntaxNode
	SetLeft(n SyntaxNode)
	Right() SyntaxNode
	SetRight(n SyntaxNode)
	AddChild(n SyntaxNode)
}

// 基础节点
type BaseNode struct {
	kind   NodeKind
	parent SyntaxNode
	left   SyntaxNode
	right  SyntaxNode
}

func (this *BaseNode) Kind() NodeKind {
	return this.kind
}

func (this *BaseNode) ChangeKind(kind NodeKind) {
	this.kind = kind
}

func (this *BaseNode) Parent() SyntaxNode {
	return this.parent
}
func (this *BaseNode) SetParent(p SyntaxNode) {
	this.parent = p
}
func (this *BaseNode) Left() SyntaxNode {
	return this.left
}
func (this *BaseNode) SetLeft(n SyntaxNode) {
	this.left = n
	n.SetParent(this)
}
func (this *BaseNode) Right() SyntaxNode {
	return this.right
}
func (this *BaseNode) SetRight(n SyntaxNode) {
	this.right = n
	n.SetParent(this)
}

func (this *BaseNode) AddChild(n SyntaxNode) {
	if this.left == nil {
		this.left = n
		n.SetParent(this)
	} else if this.right == nil {
		this.right = n
		n.SetParent(this)
	}
}

func NewSpaceNode() SyntaxNode {
	return &BaseNode{
		NodeKindSpace,
		nil,
		nil,
		nil,
	}
}
func NewAndNode() SyntaxNode {
	return &BaseNode{
		NodeKindAnd,
		nil,
		nil,
		nil,
	}
}
func NewOrNode() SyntaxNode {
	return &BaseNode{
		NodeKindOr,
		nil,
		nil,
		nil,
	}
}

type FuncNode struct {
	BaseNode
	name   string
	params []*Token
}

func NewFuncNode(name string) SyntaxNode {
	return &FuncNode{
		BaseNode{
			NodeKindFunc,
			nil,
			nil,
			nil,
		},
		name,
		make([]*Token, 0),
	}
}
func (this *FuncNode) SetFuncName(name string) {
	this.name = name
}
func (this *FuncNode) AddParam(t *Token) {
	this.params = append(this.params, t)
}
