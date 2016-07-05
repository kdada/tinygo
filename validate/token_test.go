package validate

import (
	"fmt"
	"testing"
)

func TestTokenizor(t *testing.T) {
	var expr = []byte(`(> =-623 && dd33< '1 0' &&/[a-zA-Z]\//) || isEven ( -23 , '\'  \\()&&||33 )'`)
	var tokenizor = NewTokenizor(expr)
	for {
		var token, e = tokenizor.Fetch()
		if e != nil {
			if e != EOF {
				t.Error(e)
			}
			break
		}
		fmt.Println(token.Kind)
		fmt.Println("[" + token.Value + "]")
	}
}
