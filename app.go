package tinygo

// App 应用接口
type App interface {
	// Name 返回App名称
	Name() string
	// Init 应用初始化接口
	Init() error
	// Run 应用运行接口
	Run() error
}
