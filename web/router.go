package web

import "github.com/kdada/tinygo/router"

// 创建适用于Web App的根路由
func NewRootRouter() (router.Router, error) {
	return router.NewRouter("base", "", nil)
}
