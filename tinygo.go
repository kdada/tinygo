// Package tinygo 实现了一个轻量级连接&路由框架
package tinygo

import (
	"fmt"
	"sync"

	"github.com/kdada/tinygo/config"
	"github.com/kdada/tinygo/connector"
	"github.com/kdada/tinygo/router"
)

// App 应用
type App struct {
	Connector  connector.Connector //连接器
	Root       router.Router       //根路由
	Dispatcher *Dispatcher         //调度器
}

type Context struct {
	router.BaseContext
	routerInfo map[string]string
	data       interface{}
}

// NewContext 创建上下文信息
func NewContext(segments []string, data interface{}) *Context {
	var context = new(Context)
	context.routerInfo = make(map[string]string, 1)
	context.Segs = segments
	context.data = data
	return context
}

// Value 返回路由值
func (this *Context) Value(name string) (string, bool) {
	var value, ok = this.routerInfo[name]
	return value, ok
}

// SetValue 设置路由值
func (this *Context) SetValue(name string, value string) {
	this.routerInfo[name] = value
}

// Data 返回路由上下文携带的信息
func (this *Context) Data() interface{} {
	return this.data
}

// Dispatcher 调度器,用于协调连接器和路由
type Dispatcher struct {
	Root router.Router
}

// Dispatch 分发
//  segments:用于进行分发的路径段信息
//  data:连接携带的数据
func (this *Dispatcher) Dispatch(segments []string, data interface{}) {
	var context = NewContext(segments, data)
	var executor, ok = this.Root.Match(context)
	if ok {
		var err = executor.Execute()
		if err != nil {

		}
	} else {
	}
}

// 管理器
type Manager struct {
	Apps []*App
}

// NewManager 创建管理器
func NewManager(configKind string, path string) (*Manager, error) {
	var cfg, err = config.NewConfig(configKind, path)
	if err != nil {
		return nil, err
	}
	var sections = cfg.Sections()
	var manager = new(Manager)
	manager.Apps = make([]*App, 0, 1)
	for _, s := range sections {
		var connKind, err = s.String("connector")
		if err != nil {
			panic(ErrorConfigNotCorrect.Format(s.Name(), "connector"))
		}
		connSource, err := s.String("source")
		if err != nil {
			panic(ErrorConfigNotCorrect.Format(s.Name(), "source"))
		}
		routerKind, err := s.String("router")
		if err != nil {
			panic(ErrorConfigNotCorrect.Format(s.Name(), "router"))
		}
		routerMatch, err := s.String("root")
		if err != nil {
			panic(ErrorConfigNotCorrect.Format(s.Name(), "root"))
		}
		conn, err := connector.NewConnector(connKind, connSource)
		if err != nil {
			panic(ErrorConnectorCreateFail.Format(connKind, err.Error()))
		}
		root, err := router.NewRouter(routerKind, "", routerMatch)
		if err != nil {
			panic(ErrorConnectorCreateFail.Format(routerKind, err.Error()))
		}
		manager.AddApp(conn, root)
	}
	return manager, nil
}

// NewEmptyManager 创建空管理器
func NewEmptyManager() (*Manager, error) {
	return &Manager{make([]*App, 0, 1)}, nil
}

// AddAppGroup 添加App组
func (this *Manager) AddApp(connector connector.Connector, root router.Router) *App {
	var app = new(App)
	app.Connector = connector
	app.Root = root
	app.Dispatcher = &Dispatcher{root}
	app.Connector.SetDispatcher(app.Dispatcher)
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
