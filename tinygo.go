// Package tinygo 实现了一个轻量级的Http Server框架
package tinygo

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
)

// Run 运行Http Server
func Run() {
	//加载配置
	var appFilePath, _ = exec.LookPath(os.Args[0])
	var err = loadConfig(filepath.Dir(appFilePath))
	if err != nil {
		fmt.Println(err)
		return
	}
	err = loadLayoutConfig()
	if err != nil {
		fmt.Println(err)
		return
	}
	//生成静态路由
	generateStaticRouters()
	//预编译视图
	if tinyConfig.precompile {
		compileAllViews()
	}
	if tinyConfig.session {
		//初始化Session机制
		initSession(tinyConfig.sessiontype, tinyConfig.sessionexpire)
	}
	if tinyConfig.home != "" {
		//设置首页
		SetHomePage(tinyConfig.home)
	}
	//启动
	http.HandleFunc("/", handler)
	var port = fmt.Sprintf(":%d", tinyConfig.port)
	fmt.Println("TinyGo开始监听,端口:", tinyConfig.port)
	if tinyConfig.https {
		//启动https监听
		err = http.ListenAndServeTLS(port, tinyConfig.cert, tinyConfig.pkey, nil)
	} else {
		//启动http监听
		err = http.ListenAndServe(port, nil)
	}
	if err != nil {
		fmt.Println(err)
		return
	}

}

// SafeEnvironment在安全环境中执行f方法,安全环境中出现panic不会引起进程崩溃
func SafeEnvironment(f func()) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()
	f()
}

// AsyncSafeEnvironment在一个goroutine安全环境中执行f方法,安全环境中出现panic不会引起进程崩溃
func AsyncSafeEnvironment(f func()) {
	go SafeEnvironment(f)
}
