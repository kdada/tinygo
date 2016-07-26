package validator

import (
	"fmt"
	"testing"
)

func TestParser(t *testing.T) {
	var src = `len <= 9090 && (isok() || sad() && (isNum(+222.3332432,-435345,'ae53你好24asfd\\''asdee')&&(iscc()||ssde())))||/asdasd//2323[]{}''/`
	var l = NewLexer(src)
	var p = NewParser(l)
	var err = p.Parse()
	if err != nil {
		t.Fatal(err)
	} else {
		printTree(p.Tree, "#")
	}
	var node = p.Tree
	if node.Kind() != NodeKindOr {
		t.Fatal("语法树错误")
	}
	if node.Left().Kind() != NodeKindAnd {
		t.Fatal("语法树错误")
	}
	if node.Left().Right().Kind() != NodeKindAnd {
		t.Fatal("语法树错误")
	}
	if node.Left().Right().Left().Kind() != NodeKindOr {
		t.Fatal("语法树错误")
	}
	if node.Left().Right().Right().Kind() != NodeKindAnd {
		t.Fatal("语法树错误")
	}
	if node.Left().Right().Right().Right().Kind() != NodeKindOr {
		t.Fatal("语法树错误")
	}
}

type Stringer interface {
	String() string
}

func (this *BaseNode) String() string {
	return fmt.Sprintf("%d", int(this.kind))
}
func (this *FuncNode) String() string {
	var f = this.name
	for _, v := range this.params {
		f += "  " + fmt.Sprint(v.Value)
	}
	return f
}

func printTree(root SyntaxNode, pre string) {
	if root == nil {
		return
	}
	fmt.Print(pre + ":")
	fmt.Println(root.(Stringer).String())
	if root.Left() != nil {
		printTree(root.Left(), pre+"L-")
	}
	if root.Right() != nil {
		printTree(root.Right(), pre+"R-")
	}
}
