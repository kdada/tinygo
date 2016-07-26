package validator

import "reflect"

// 字符串验证器
type StringValidator struct {
	Tree SyntaxNode
}

// 创建字符串验证器
func NewStringValidator(source string) (Validator, error) {
	var p = NewParser(NewLexer(source))
	var err = p.Parse()
	if err != nil {
		return nil, err
	}
	var root = NewSpaceNode()
	root.SetLeft(p.Tree)
	err = translate(root.Left())
	if err != nil {
		return nil, err
	}
	var t = root.Left()
	t.SetParent(nil)
	root.SetLeft(nil)
	return &StringValidator{t}, nil
}

// translate 将所有func和regfunc节点全部转换为ValidatorFunc
func translate(node SyntaxNode) error {
	switch node.Kind() {
	case NodeKindFunc, NodeKindRegFunc:
		var f ValidatorFunc
		var err error
		var params []interface{}
		var fnode = node.(*FuncNode)
		if fnode.Kind() == NodeKindFunc {
			//处理参数信息
			var name = fnode.name + sep
			for _, p := range fnode.params {
				var k = reflect.Invalid
				switch p.Kind {
				case TokenKindInteger:
					k = reflect.Int64
				case TokenKindFloat:
					k = reflect.Float64
				case TokenKindString:
					k = reflect.String
				}
				var e, ok = CheckType(k)
				if ok {
					name += e
				} else {
					return ErrorIllegalParam.Format(p.Kind).Error()
				}
			}
			var vf, ok = funcs[name]
			if ok {
				f = vf
			} else {
				err = ErrorUnmatchedFunc.Format(name).Error()
			}
		} else {
			f, err = NewRegFunc(fnode.name)
			params = []interface{}{}
		}
		if err != nil {
			return err
		}
		var newNode = NewExecutableFuncNode(f, params)
		if node.Parent().Left() == node {
			node.Parent().SetLeft(newNode)
		} else {
			node.Parent().SetRight(newNode)
		}
		return nil
	case NodeKindAnd, NodeKindOr:
		var err = translate(node.Left())
		if err != nil {
			return err
		}
		err = translate(node.Right())
		return err
	}
	return ErrorIllegalNode.Format(node.Kind()).Error()
}

// 验证
func (this *StringValidator) Validate(str string) bool {
	return this.validate(this.Tree, str)
}

// 递归验证
func (this *StringValidator) validate(node SyntaxNode, str string) bool {
	if node.Kind() == NodeKindExecutor {

	}
	switch node.Kind() {
	case NodeKindExecutor:
		return node.(*ExecutableFuncNode).Execute(str)
	case NodeKindAnd:
		var result = this.validate(node.Left(), str)
		if result == false {
			return false
		}
		return this.validate(node.Right(), str)
	case NodeKindOr:
		var result = this.validate(node.Left(), str)
		if result == true {
			return true
		}
		return this.validate(node.Right(), str)
	}
	panic(ErrorIllegalNode.Format(node.Kind()))
}
