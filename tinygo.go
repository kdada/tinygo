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

// HttpNotFound 返回页面不存在(404)错误
func HttpNotFound(w http.ResponseWriter, r *http.Request) {
	if tinyConfig.pageerr != "" {
		w.WriteHeader(404)
		ParseTemplate(w, r, tinyConfig.pageerr, nil)
	} else {
		http.NotFound(w, r)
	}
}

// Redirect 302重定向
func Redirect(w http.ResponseWriter, r *http.Request, url string) {
	http.Redirect(w, r, url, 302)
}

// RedirectPermanently 301重定向
func RedirectPermanently(w http.ResponseWriter, r *http.Request, url string) {
	http.Redirect(w, r, url, 301)
}
