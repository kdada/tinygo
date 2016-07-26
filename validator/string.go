package validator

import "reflect"

// 字符串验证器
//  1.使用()改变优先级
//  2.使用&&和||连接函数或表达式
//  3.函数包括以下几种
//    (1)普通函数:IsOK IsOK() BigThan(1234)  Contain('abc')
//    (2)名称中包含关系运算符:>=10 Len==11 Complex<(12,22)
//    (3)正则表达式:/[a-z]+?/
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
			params = make([]interface{}, len(fnode.params))
			for i, p := range fnode.params {
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
					params[i] = p.Value
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
