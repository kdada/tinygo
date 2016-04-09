package tinygo

import (
	"encoding/json"
	"io/ioutil"
	"path"
	"path/filepath"
	"strings"

	"github.com/kdada/tinygo/config"
)

// NormalizePath 规范化路径
// 规范化后统一使用/分隔符,所有绝对路径全部变为相对路径
//
// 例如:
//  \test\test.html ==> test/test.html
func NormalizePath(path string) string {
	return strings.Trim(filepath.ToSlash(path), "/")
}

// 基本配置
var tinyConfig = struct {
	app           string   //应用名称
	path          string   //当前程序启动目录(无需从文件读取)
	mode          string   //启动模式,可以为debug或release
	https         bool     //是否启用https,可选,默认为false
	port          uint16   //监听端口,可选,默认为80，https为true则默认为443
	cert          string   //证书(PEM)路径,如果启用了https则必填
	pkey          string   //私钥(PEM)路径,如果启用了https则必填
	home          string   //首页地址
	session       bool     //是否启用session
	sessiontype   string   //session类型,参考tinygo/session,默认为memory
	sessionexpire int64    //session过期时间,单位为秒
	csrf          bool     //是否启用csrf
	csrfexpire    int64    //csrf token过期时间
	static        []string //静态文件目录,默认为"content"
	view          string   //视图文件目录,默认为"views"
	pageerr       string   //默认错误页面路径,默认为空
	precompile    bool     //是否预编译页面路径,默认为false
	api           string   //使用Api返回的数据的解析格式,默认为auto
	//自动设置项
	sessionName string //session对应的Cookie名称 app+DefaultSessionCookieName
	csrfName    string //csrf对应的Cookie名称 app+DefaultCSRFCookieName
}{}

// loadConfig 加载配置
//  envPath:当前程序启动目录
func loadConfig(envPath string) error {
	tinyConfig.path = envPath
	var configPath = filepath.Join(tinyConfig.path, DefaultConfigPath)
	var cfg, err = config.NewConfig(config.ConfigTypeIni, configPath)
	if err != nil {
		//配置文件加载出错
		return err
	} else {
		//读取配置文件
		var global = cfg.GlobalSection()
		var err error
		//app
		tinyConfig.app, err = global.String("app")
		if err != nil {
			tinyConfig.app = "app"
		}
		//设置name
		tinyConfig.sessionName = tinyConfig.app + DefaultSessionCookieName
		tinyConfig.csrfName = tinyConfig.app + DefaultCSRFCookieName

		//mode
		tinyConfig.mode, err = global.String("mode")
		if err != nil {
			tinyConfig.mode = "debug"
		}
		//https
		tinyConfig.https, err = global.Bool("https")
		if err != nil {
			tinyConfig.https = false
		}
		//port
		port, err := global.Int("port")
		if err != nil {
			if tinyConfig.https {
				port = 443
			} else {
				port = 80
			}
		}
		tinyConfig.port = uint16(port)
		if tinyConfig.https {
			//cert
			tinyConfig.cert, _ = global.String("cert")
			//key
			tinyConfig.pkey, _ = global.String("pkey")
		}
		//home
		tinyConfig.home, err = global.String("home")
		if err != nil {
			tinyConfig.home = "/home/index"
		}
		//session
		tinyConfig.session, err = global.Bool("session")
		if err != nil {
			tinyConfig.session = false
		}
		//sessiontype
		tinyConfig.sessiontype, err = global.String("sessiontype")
		if err != nil {
			tinyConfig.sessiontype = "memory"
		}
		//sessionexpire
		tinyConfig.sessionexpire, err = global.Int("sessionexpire")
		if err != nil {
			tinyConfig.sessionexpire = 3600
		}
		//csrf
		tinyConfig.csrf, err = global.Bool("csrf")
		if err != nil {
			tinyConfig.csrf = false
		}
		//csrfexpire
		tinyConfig.csrfexpire, err = global.Int("csrfexpire")
		if err != nil {
			tinyConfig.csrfexpire = 3600
		}
		//static
		paths, err := global.String("static")
		if err != nil {
			tinyConfig.static = []string{"content"}
		} else {
			tinyConfig.static = strings.Split(paths, ";")
			//路径规范化
			//统一规范为相对路径
			// \content\css ==> content/css
			for i, p := range tinyConfig.static {
				tinyConfig.static[i] = NormalizePath(p)
			}
		}
		//view
		tinyConfig.view, err = global.String("view")
		if err != nil {
			tinyConfig.view = "views"
		}
		//pageerr
		tinyConfig.pageerr, _ = global.String("pageerr")
		//precompile
		tinyConfig.precompile, err = global.Bool("precompile")
		if err != nil {
			tinyConfig.precompile = false
		}
		//api
		tinyConfig.api, err = global.String("api")
		if err != nil {
			tinyConfig.api = "auto"
		}
		return nil
	}
}

// 判断当前是否是发布模式
// 由配置文件的mode决定
func IsRelease() bool {
	return tinyConfig.mode == "release"
}

// getViewFilePath 返回视图文件路径的绝对路径
//  viewFilePath:相对于视图目录的文件路径
//  layout/layout.html ==> {{tinyConfig.view}}/layout/layout.html
func getViewFilePath(viewFilePath string) string {
	return filepath.Join(tinyConfig.view, viewFilePath)
}

// generateViewFilePath 返回视图文件相对tinyConfig.view的路径
//  viewFilePath:视图文件绝对路径
//  {{tinyConfig.view}}/layout/layout.html ==> layout/layout.html
func generateViewFilePath(viewFilePath string) string {
	viewFilePath, _ = filepath.Rel(tinyConfig.view, viewFilePath)
	return viewFilePath
}

// 视图布局配置(json格式)
// 为了确保顺利解析,即使在不同目录中
// layout文件不得与任何视图文件重名,layout文件之间也不允许重名
var layoutConfig = struct {
	LayoutMap     map[string]string //为布局文件定义别名,map[布局文件别名]布局文件
	DefaultLayout string            //默认布局文件别名
	LayoutSpec    map[string]string //特别指定目录或文件对应的布局文件,map[目录或文件]布局文件别名
}{
	LayoutMap:     map[string]string{},
	DefaultLayout: "",
	LayoutSpec:    map[string]string{},
}

// loadLayoutConfig 加载布局配置
func loadLayoutConfig() error {
	var configPath = getViewFilePath(DefaultLayoutConfigFileName)
	var content, err = ioutil.ReadFile(configPath)
	if err != nil {
		Log("Tinygo use default layout config")
		content = []byte(`{"LayoutMap":{},"DefaultLayout":"","LayoutSpec":{}}`)
	}
	err = json.Unmarshal(content, &layoutConfig)
	if err != nil {
		return err
	} else {
		//路径规范化
		//统一规范为相对路径
		// \layout\layout.html ==> layout/layout.html
		// \layout\ ==> layout
		for k, v := range layoutConfig.LayoutMap {
			layoutConfig.LayoutMap[k] = NormalizePath(v)
		}
		for k, v := range layoutConfig.LayoutSpec {
			delete(layoutConfig.LayoutSpec, k)
			layoutConfig.LayoutSpec[NormalizePath(k)] = v
		}
		return nil
	}
}

// isLayoutFile 判断文件是否是布局文件
// filePath:相对于tinyConfig.view的文件路径
func isLayoutFile(filePath string) bool {
	filePath = NormalizePath(filePath)
	for _, layoutPath := range layoutConfig.LayoutMap {
		if filePath == layoutPath {
			return true
		}
	}
	return false
}

// getLayout 获取指定文件的Layout文件
//  filePath:相对于tinyConfig.view的文件路径
func getLayoutFile(filePath string) (string, bool) {
	filePath = NormalizePath(filePath)
	var layout, ok = layoutConfig.LayoutSpec[filePath]
	if !ok {
		var dir = path.Dir(filePath)
		layout, ok = layoutConfig.LayoutSpec[dir]
		if !ok {
			layout = layoutConfig.DefaultLayout
		}
	}
	layout, ok = layoutConfig.LayoutMap[layout]
	if filePath == layout {
		//避免循环布局
		return "", false
	}
	return layout, ok
}
