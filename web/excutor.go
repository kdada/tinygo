package web

import "github.com/kdada/tinygo/router"

// MethodExecutor
type MethodExecutor struct {
	router.BaseRouterExecutor
	Method func(end router.Router, context router.RouterContext) error
}

// Excute 执行
func (this *MethodExecutor) Execute() error {
	return this.Method(this.End, this.Context)
}
