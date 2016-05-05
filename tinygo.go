// Package tinygo 实现了一个轻量级Http框架
package tinygo

import "github.com/kdada/tinygo/app"

var Manager, _ = app.NewManager()

func Run() error {
	return nil
}
