// Package tinygo 实现了一个轻量级Http框架
package tinygo

import (
	"strconv"

	"github.com/kdada/tinygo/app"
	"github.com/kdada/tinygo/connector"
	"github.com/kdada/tinygo/router"
)

var Manager, _ = app.NewManager()

func Run(appDir string, configFile string, root router.Router) error {
	//读取配置文件
	var config, err = ReadHttpConfig(appDir, configFile)
	if err != nil {
		return err
	}
	var conn connector.Connector
	if config.Https {
		panic("https not defined")
	} else {
		conn, err = connector.NewConnector("http", ":"+strconv.Itoa(config.Port))
		if err != nil {
			return err
		}
	}
	var dispatcher = NewHttpProcessor(root, config)
	var a = app.NewApp(conn, dispatcher)
	Manager.AddApp(a)
	Manager.Run()
	return nil
}
