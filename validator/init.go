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
	RegisterFunc("len<", lenLtI)
	RegisterFunc("len<=", lenLeI)
	RegisterFunc("len>", lenGtI)
	RegisterFunc("len>=", lenGeI)
	RegisterFunc("len==", lenEqI)
	RegisterFunc("len!=", lenNeI)
	RegisterFunc("clen<", clenLtI)
	RegisterFunc("clen<=", clenLeI)
	RegisterFunc("clen>", clenGtI)
	RegisterFunc("clen>=", clenGeI)
	RegisterFunc("clen==", clenEqI)
	RegisterFunc("clen!=", clenNeI)

	Register("string", NewStringValidator)
}
