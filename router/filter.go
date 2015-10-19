package router

// 路由过滤器
type RouterFilter interface {
	// Filter 过滤该请求
	// return:返回true表示继续处理,否则终止路由过程,后续的过滤器也不会执行
	Filter(context RouterContext) bool
}
