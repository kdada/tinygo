package web

import (
	"strconv"

	"github.com/kdada/tinygo/connector"
	"github.com/kdada/tinygo/router"
)

// Web应用
type WebApp struct {
	conn      connector.Connector
	processor *HttpProcessor
}

// NewWebApp 创建WebApp
//  appDir:配置文件路径
//  configFile:配置文件名称
//  root:Web使用的根路由,根路由必须是匹配空字符串的路由
func NewWebApp(appDir string, configFile string, root router.Router) (*WebApp, error) {
	//读取配置文件
	var config, err = ReadHttpConfig(appDir, configFile)
	if err != nil {
		return nil, err
	}
	var conn connector.Connector
	if config.Https {
		panic("https not defined")
	} else {
		conn, err = connector.NewConnector("http", ":"+strconv.Itoa(config.Port))
		if err != nil {
			return nil, err
		}
	}
	var processor = NewHttpProcessor(root, config)
	conn.SetDispatcher(processor)
	return &WebApp{conn, processor}, nil
}

// Name 应用名称
func (this *WebApp) Name() string {
	return this.processor.Config.App
}

// Init 应用初始化接口
func (this *WebApp) Init() error {
	return this.conn.Init()
}

// Run 应用运行接口
func (this *WebApp) Run() error {
	return this.conn.Run()
}
