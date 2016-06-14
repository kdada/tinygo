package web

type Controller interface {
	// Start 在进入控制器时执行
	Start()
	// End 在退出控制器时执行
	End()
}
