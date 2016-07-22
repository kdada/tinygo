package validator

import (
	"fmt"
	"testing"
)

func TestParser(t *testing.T) {
	var src = `<=9090 && (isok()||isNum(+222.33,-435345,'ae5324asfd\'asdee'))&&/asdasd\/2323[]{}''/`
	var l = NewLexer(src)
	var p = NewParser(l)
	p.Parse()
	printTree(p.Tree, "")

}

type Stringer interface {
	String() string
}

func (this *BaseNode) String() string {
	return string(this.kind)
}
func (this *FuncNode) String() string {
	var f = this.name + "  "
	for _, v := range this.params {
		f += fmt.Sprint(v)
	}
	return f
}

func printTree(root SyntaxNode, pre string) {
	if root == nil {
		return
	}
	fmt.Print(pre + string(root.Kind()) + ":")
	fmt.Println(root.(Stringer).String())
	if root.Left() != nil {
		printTree(root.Left(), pre+"L-")
	}
	if root.Right() != nil {
		printTree(root.Right(), pre+"R-")
	}
}
