package tinygo

import (
	"fmt"
	"sync"

	"github.com/kdada/tinygo/connector"
	"github.com/kdada/tinygo/router"
)

// App
type App struct {
	Connector connector.Connector //连接器
	Root      router.Router       //根路由
}

// 管理器
type Manager struct {
	Apps []*App
}

// AddAppGroup 添加App组
func (this *Manager) AddApp(connector connector.Connector, root router.Router) *App {
	var app = new(App)
	app.Connector = connector
	app.Root = root
	this.Apps = append(this.Apps, app)
	return app
}

// Run 运行App
func (this *Manager) Run() error {
	var err error
	var w sync.WaitGroup
	w.Add(len(this.Apps))
	for _, app := range this.Apps {
		err = app.Connector.Init()
		if err != nil {
			return err
		}
		go func() {
			defer func() {
				if err := recover(); err != nil {
					fmt.Println(err)
				}
				w.Done()
			}()
			app.Connector.Run()
		}()
	}
	w.Wait()
	return nil
}
