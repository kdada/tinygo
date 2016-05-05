// Package app 实现了一个应用管理器
package app

import (
	"github.com/kdada/tinygo/connector"
	"github.com/kdada/tinygo/router"
)

// 调度器用于关联连接器和路由
type Dispatcher interface {
	// Dispatch 进行调度
	//  segments:用于进行路由的路径段
	//  data:连接器传递的数据
	Dispatch(segments []string, data interface{})
}

// App 应用
//  引用(->)关系:
//  连接器->调度器->根路由
type App struct {
	Connector  connector.Connector //连接器
	Root       router.Router       //根路由
	Dispatcher Dispatcher          //调度器
}

// NewApp 创建新的App
func NewApp(connector connector.Connector, root router.Router, dispatcher Dispatcher) *App {
	var app = &App{connector, root, dispatcher}
	app.Connector.SetDispatcher(dispatcher)
	return app
}
