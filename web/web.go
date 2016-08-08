// package web 实现了处理http请求的基本工具
package web

import (
	"strconv"

	"github.com/kdada/tinygo/connector"
	"github.com/kdada/tinygo/router"
)

// Web应用
type WebApp struct {
	Conn      connector.Connector //链接
	Processor *HttpProcessor      //处理器
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
		conn, err = connector.NewConnector("https", ":"+strconv.Itoa(config.Port)+";Cert="+config.Cert+";Key="+config.PrivateKey)
	} else {
		conn, err = connector.NewConnector("http", ":"+strconv.Itoa(config.Port))
	}
	if err != nil {
		return nil, err
	}
	var processor, e = NewHttpProcessor(root, config)
	if e != nil {
		return nil, e
	}
	conn.SetDispatcher(processor)
	return &WebApp{conn, processor}, nil
}

// Name 应用名称
func (this *WebApp) Name() string {
	return this.Processor.Config.App
}

// Init 应用初始化接口
func (this *WebApp) Init() error {
	return this.Conn.Init()
}

// Run 应用运行接口
func (this *WebApp) Run() error {
	return this.Conn.Run()
}
