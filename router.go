package tinygo

import (
	"github.com/kdada/tinygo/info"
	"github.com/kdada/tinygo/router"
)

//根路由,其他路由应该作为根路由的子路由
var RootRouter = router.NewRootRouter()

// 控制器方法路由信息
type RouterInfo struct {
	MethodName string          //方法名
	HttpMethod info.HttpMethod //http方法
	Extensions []string        //url扩展
}

func (this *RouterInfo) Info() (string, string, []string) {
	return this.MethodName, string(this.HttpMethod), this.Extensions
}
