package validate

import (
	"fmt"
	"testing"
)

func TestParser(t *testing.T) {
	var expr = []byte(`isEven(-123,'asdasd\ddd'||(>=-623&&||len<10&&/[0-9]*/)`)
	var tokenizor = NewTokenizor(expr)
	var root, err = Parse(tokenizor)
	if err != nil && err != EOF {
		t.Error(err)
	}
	Print(root, "")
}

func Print(node *ScopeNode, level string) {
	if node == nil {
		return
	}
	fmt.Print(level + string(node.Kind) + ":")
	fmt.Println(node.Value)
	if node.HasLeft() {
		Print(node.Left, level+"L-")
	}
	if node.HasRight() {
		Print(node.Right, level+"R-")
	}
}
