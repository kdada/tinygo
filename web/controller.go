package web

type Controller interface {
	// Start 在进入控制器时执行
	Start()
	// Finish 在退出控制器时执行
	Finish()
}
