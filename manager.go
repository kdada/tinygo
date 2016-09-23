package tinygo

import "sync"

// 管理器
type Manager struct {
	Apps  []App
	Event ManagerEvent
}

// NewManager 创建App管理器
func NewManager() (*Manager, error) {
	return &Manager{make([]App, 0, 1), &DefaultManagerEvent{}}, nil
}

// AddApp 添加App
func (this *Manager) AddApp(app App) {
	this.Apps = append(this.Apps, app)
}

// Run 运行App
func (this *Manager) Run() {
	var err error
	if this.Event != nil {
		this.Event.Started(this)
		defer this.Event.Ended(this)
	}
	var w sync.WaitGroup
	w.Add(len(this.Apps))
	for _, app := range this.Apps {
		if this.Event != nil {
			this.Event.BeforeRunning(this, app)
		}
		err = app.Init()
		if err != nil {
			this.Event.InitFailed(this, app, err)
			continue
		}
		go func(app App) {
			defer func() {
				if err := recover(); err != nil {
					if this.Event != nil {
						this.Event.RunPaniced(this, app, err)
					}
				}
				w.Done()
			}()
			var err = app.Run()
			if err != nil {
				if this.Event != nil {
					this.Event.RunFailed(this, app, err)
				}
			} else {
				if this.Event != nil {
					this.Event.RunFinished(this, app)
				}
			}
		}(app)
		if this.Event != nil {
			this.Event.AferRunning(this, app)
		}
	}
	w.Wait()
}
