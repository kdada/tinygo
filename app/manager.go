package app

import (
	"fmt"
	"sync"
)

// 管理器事件
type ManagerEvent interface {
	// Started 在Manager运行的时候触发
	Started(manager *Manager)
	// BeforeAppRunning 在App运行前触发
	BeforeAppRunning(manager *Manager, app *App)
	// AferAppRunning 在App运行后触发
	AferAppRunning(manager *Manager, app *App)
	// Ended 在Manager停止运行时触发
	Ended(manager *Manager)
}

// 管理器
type Manager struct {
	Apps  []*App
	Event ManagerEvent
}

// NewManager 创建App管理器
func NewManager() (*Manager, error) {
	return &Manager{make([]*App, 0, 1), nil}, nil
}

// AddApp 添加App
func (this *Manager) AddApp(app *App) {
	this.Apps = append(this.Apps, app)
}

// Run 运行App
func (this *Manager) Run() error {
	var err error
	if this.Event != nil {
		this.Event.Started(this)
		defer this.Event.Ended(this)
	}
	var w sync.WaitGroup
	w.Add(len(this.Apps))
	for _, app := range this.Apps {
		if this.Event != nil {
			this.Event.BeforeAppRunning(this, app)
		}
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
		if this.Event != nil {
			this.Event.AferAppRunning(this, app)
		}
	}
	w.Wait()

	return nil
}
