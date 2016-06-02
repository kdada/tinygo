// Package tinygo 实现一个组合式应用管理器
package tinygo

// tinygo App管理器
var manager, _ = NewManager()

// SetEvent 设置App管理器事件
func SetEvent(event ManagerEvent) {
	manager.Event = event
}

// AddApp 添加App
func AddApp(app App) {
	manager.AddApp(app)
}

// AddApps 批量添加App
func AddApps(apps ...App) {
	for _, app := range apps {
		AddApp(app)
	}
}

// Run 运行
func Run() {
	manager.Run()
}
