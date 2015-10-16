package tinygo

import (
	"encoding/json"
	"io/ioutil"
	"path"
	"path/filepath"
	"strings"

	"github.com/kdada/tinygo/config"
	"github.com/kdada/tinygo/info"
)

// NormalizePath 规范化路径
// 规范化后统一使用/分隔符,所有绝对路径全部变为相对路径
// 例如: \test\test.html ==> test/test.html
func NormalizePath(path string) string {
	return strings.TrimLeft(filepath.ToSlash(path), "/")
}

// 基本配置(ini格式)
var tinyConfig = struct {
	path       string   //当前程序启动目录
	https      bool     //是否启用https,可选,默认为false
	port       uint16   //监听端口,可选,默认为80，https为true则默认为443
	cert       string   //证书(PEM)路径,如果启用了https则必填
	key        string   //私钥(PEM)路径,如果启用了https则必填
	session    string   //session类型,参考tinygo/session,默认为memory
	static     []string //静态文件目录,默认为"content"
	view       string   //视图文件目录,默认为"views"
	pageerr    string   //默认错误页面路径,默认为空
	precompile bool     //是否预编译页面路径,默认为false
	api        string   //使用Api返回的数据的解析格式,默认为auto
}{}

// loadConfig 加载配置
// envPath:当前程序启动目录
func loadConfig(envPath string) error {
	tinyConfig.path = envPath
	var configPath = filepath.Join(tinyConfig.path, info.DefaultConfigPath)
	var cfg, err = config.NewConfig(config.ConfigTypeIni, configPath)
	if err != nil {
		//配置文件加载出错
		return err
	} else {
		//读取配置文件
		var global = cfg.GlobalSection()
		var err error
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
			tinyConfig.key, _ = global.String("key")
		}
		//session
		tinyConfig.session, err = global.String("session")
		if err != nil {
			tinyConfig.session = "memory"
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

// getViewFilePath 返回视图文件路径的绝对路径
// viewFilePath:相对于视图目录的文件路径
// layout/layout.html ==> {{tinyConfig.path}}/{{tinyConfig.view}}/layout/layout.html
func getViewFilePath(viewFilePath string) string {
	return filepath.Join(tinyConfig.path, tinyConfig.view, viewFilePath)
}

// generateViewFilePath 返回视图文件相对tinyConfig.view的路径
// viewFilePath:视图文件绝对路径
// {{tinyConfig.path}}/{{tinyConfig.view}}/layout/layout.html ==> layout/layout.html
func generateViewFilePath(viewFilePath string) string {
	viewFilePath, _ = filepath.Rel(filepath.Join(tinyConfig.path, tinyConfig.view), viewFilePath)
	return viewFilePath
}

// 视图布局配置(json格式)
// 为了确保顺利解析,即使在不同目录中
// layout文件不得与任何视图文件重名,layout文件之间也不允许重名
var layoutConfig = struct {
	LayoutMap     map[string]string //为布局文件定义别名
	DefaultLayout string            //默认布局文件别名
	LayoutSpec    map[string]string //特别指定目录或文件对应的布局文件,map[目录或文件]布局文件别名
}{
	LayoutMap:     map[string]string{},
	DefaultLayout: "",
	LayoutSpec:    map[string]string{},
}

// loadLayoutConfig 加载布局配置
func loadLayoutConfig() error {
	var configPath = getViewFilePath(info.DefaultLayoutConfigFileName)
	var content, err = ioutil.ReadFile(configPath)
	if err == nil {
		err = json.Unmarshal(content, &layoutConfig)
	}
	if err != nil {
		return err
	} else {
		//路径规范化
		//统一规范为相对路径
		// \layout\layout.html ==> layout/layout.html
		for k, v := range layoutConfig.LayoutMap {
			layoutConfig.LayoutMap[k] = NormalizePath(v)
		}
		for k, v := range layoutConfig.LayoutSpec {
			layoutConfig.LayoutSpec[k] = NormalizePath(v)
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
// filePath:相对于tinyConfig.view的文件路径
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
