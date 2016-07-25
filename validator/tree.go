package validator

// 节点类型
type NodeKind string

const (
	NodeKindSpace NodeKind = "NodeKindSpace" //空间节点
	NodeKindAnd   NodeKind = "NodeKindAnd"   //逻辑与节点
	NodeKindOr    NodeKind = "NodeKindOr"    //逻辑或节点
	NodeKindFunc  NodeKind = "NodeKindFunc"  //函数节点
)

// 语法节点接口
type SyntaxNode interface {
	// Kind 语法节点类型
	Kind() NodeKind
	// ChangeKind 改变语法节点类型
	ChangeKind(kind NodeKind)
	// Parent 获取父节点
	Parent() SyntaxNode
	// SetParent 设置父节点
	SetParent(n SyntaxNode)
	// Left 返回左节点
	Left() SyntaxNode
	// SetLeft 设置左节点
	SetLeft(n SyntaxNode)
	// Right 返回右节点
	Right() SyntaxNode
	// SetRight 设置右节点
	SetRight(n SyntaxNode)
	// AddChild 添加子节点
	AddChild(n SyntaxNode)
}

// 基础节点
type BaseNode struct {
	kind   NodeKind   //节点类型
	parent SyntaxNode //父节点
	left   SyntaxNode //左节点
	right  SyntaxNode //右节点
}

// Kind 语法节点类型
func (this *BaseNode) Kind() NodeKind {
	return this.kind
}

// ChangeKind 改变语法节点类型
func (this *BaseNode) ChangeKind(kind NodeKind) {
	this.kind = kind
}

// Parent 获取父节点
func (this *BaseNode) Parent() SyntaxNode {
	return this.parent
}

// SetParent 设置父节点
func (this *BaseNode) SetParent(p SyntaxNode) {
	this.parent = p
}

// Left 返回左节点
func (this *BaseNode) Left() SyntaxNode {
	return this.left
}

// SetLeft 设置左节点
func (this *BaseNode) SetLeft(n SyntaxNode) {
	this.left = n
	n.SetParent(this)
}

// Right 返回右节点
func (this *BaseNode) Right() SyntaxNode {
	return this.right
}

// SetRight 设置右节点
func (this *BaseNode) SetRight(n SyntaxNode) {
	this.right = n
	n.SetParent(this)
}

// AddChild 添加子节点
func (this *BaseNode) AddChild(n SyntaxNode) {
	if this.left == nil {
		this.left = n
		n.SetParent(this)
	} else if this.right == nil {
		this.right = n
		n.SetParent(this)
	}
}

// NewSpaceNode 创建空间节点
func NewSpaceNode() SyntaxNode {
	return &BaseNode{
		NodeKindSpace,
		nil,
		nil,
		nil,
	}
}

// NewAndNode 创建逻辑与节点
func NewAndNode() SyntaxNode {
	return &BaseNode{
		NodeKindAnd,
		nil,
		nil,
		nil,
	}
}

// NewOrNode 创建逻辑或节点
func NewOrNode() SyntaxNode {
	return &BaseNode{
		NodeKindOr,
		nil,
		nil,
		nil,
	}
}

// 函数节点
type FuncNode struct {
	BaseNode
	name   string
	params []*Token
}

// NewFuncNode 创建函数节点
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

// SetFuncName 设置函数名称
func (this *FuncNode) SetFuncName(name string) {
	this.name = name
}

// AddParam 添加参数信息
func (this *FuncNode) AddParam(t *Token) {
	this.params = append(this.params, t)
}
