package validator

func init() {
	// 注册基础验证器方法
	RegisterFunc("<", ltI)
	RegisterFunc("<=", leI)
	RegisterFunc(">", gtI)
	RegisterFunc(">=", geI)
	RegisterFunc("==", eqI)
	RegisterFunc("!=", neI)
	RegisterFunc("<", ltF)
	RegisterFunc("<=", leF)
	RegisterFunc(">", gtF)
	RegisterFunc(">=", geF)
	RegisterFunc("==", eqF)
	RegisterFunc("!=", neF)
	RegisterFunc("==", eqS)
	RegisterFunc("!=", neS)

	Register("string", NewStringValidator)
}
