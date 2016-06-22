package web

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/kdada/tinygo/config"
)

// Http配置
type HttpConfig struct {
	Root              string      //应用根目录
	App               string      //应用名称
	Mode              string      //启动模式,可以为debug或release
	Https             bool        //是否启用https,可选,默认为false
	Port              int         //监听端口,可选,默认为80，https为true则默认为443
	Cert              string      //证书(PEM)路径,如果启用了https则必填
	PrivateKey        string      //私钥(PEM)路径,如果启用了https则必填
	Home              string      //首页地址
	Session           bool        //是否启用session
	SessionType       string      //session类型,参考tinygo/session,默认为memory
	SessionSource     string      //session源,参考tinygo/session,默认为空
	SessionExpire     int         //session过期时间,单位为秒
	CSRF              bool        //是否启用csrf
	CSRFType          string      //session类型,参考tinygo/session,默认为memory
	CSRFSource        string      //session源,参考tinygo/session,默认为空
	CSRFExpire        int         //csrf token过期时间,单位为秒
	Static            []string    //静态文件目录,默认为"content",路径相对于应用根目录
	View              string      //视图文件目录,默认为"views"
	Precompile        bool        //是否预编译视图,默认为false
	Api               string      //使用Api返回的数据的解析格式,默认为auto(其他设置包括json,xml)
	Favicon           string      //网站图标路径
	Robots            string      //爬虫协议文件路径
	Log               bool        //是否启用日志
	LogType           string      //日志类型,可以为console或file
	LogPath           string      //日志路径,日志类型为file的时候需要
	LogAsync          bool        //异步日志,默认为false
	LayoutConfigPath  string      //布局配置文件名
	TemplateExt       string      //视图文件扩展名
	TemplateName      string      //模板文件内模板名,用于返回部分视图时使用
	SessionCookieName string      //Session Cookie名
	CSRFCookieName    string      //CSRF Cookie 名
	CSRFTokenName     string      //CSRF 表单名
	MaxRequestMemory  int         //单次请求最大占用内存大小,默认32 MB
	ViewConfig        *ViewConfig //视图配置
}

// NewHttpConfig 创建默认的Http配置
func NewHttpConfig() *HttpConfig {
	// Http配置
	return &HttpConfig{
		App:               "app",
		Mode:              "debug",
		Https:             false,
		Port:              80,
		Cert:              "",
		PrivateKey:        "",
		Home:              "",
		Session:           true,
		SessionType:       "memory",
		SessionSource:     "",
		SessionExpire:     1800,
		CSRF:              false,
		CSRFType:          "memory",
		CSRFSource:        "",
		CSRFExpire:        300,
		Static:            []string{"content"},
		View:              "views",
		Precompile:        false,
		Api:               "json",
		Favicon:           "favicon.icon",
		Robots:            "robots.txt",
		Log:               true,
		LogType:           "console",
		LogPath:           "",
		LogAsync:          false,
		LayoutConfigPath:  "views/layout.json",
		TemplateExt:       "html",
		TemplateName:      "Content",
		SessionCookieName: "ssid",
		CSRFCookieName:    "xid",
		CSRFTokenName:     "csrf",
		MaxRequestMemory:  32 << 20,
		ViewConfig:        &ViewConfig{make(map[string]string), "", make(map[string]string)},
	}
}

// ReadHttpConfig 读取配置文件
func ReadHttpConfig(appDir string, configPath string) (*HttpConfig, error) {
	var cfg, err = config.NewConfig("ini", filepath.Join(appDir, configPath))
	var httpCfg = NewHttpConfig()
	if err != nil {
		//配置文件加载出错
		return nil, err
	}
	httpCfg.Root = appDir
	//读取配置文件
	var global = cfg.GlobalSection()
	var strValue string
	var intValue int
	var boolValue bool

	strValue, err = global.String("App")
	if err == nil {
		httpCfg.App = strValue
	}
	strValue, err = global.String("Mode")
	if err == nil {
		httpCfg.Mode = strValue
	}
	boolValue, err = global.Bool("Https")
	if err == nil {
		httpCfg.Https = boolValue
	}
	intValue, err = global.Int("Port")
	if err == nil {
		httpCfg.Port = intValue
	} else {
		if httpCfg.Https {
			httpCfg.Port = 443
		}
	}
	strValue, err = global.String("Cert")
	if err == nil {
		httpCfg.Cert = strValue
	}
	strValue, err = global.String("PrivateKey")
	if err == nil {
		httpCfg.PrivateKey = strValue
	}
	strValue, err = global.String("Home")
	if err == nil {
		httpCfg.Home = strValue
	}
	boolValue, err = global.Bool("Session")
	if err == nil {
		httpCfg.Session = boolValue
	}
	strValue, err = global.String("SessionType")
	if err == nil {
		httpCfg.SessionType = strValue
	}
	strValue, err = global.String("SessionSource")
	if err == nil {
		httpCfg.SessionSource = strValue
	}
	intValue, err = global.Int("SessionExpire")
	if err == nil {
		httpCfg.SessionExpire = intValue
	}
	boolValue, err = global.Bool("CSRF")
	if err == nil {
		httpCfg.CSRF = boolValue
	}
	strValue, err = global.String("CSRFType")
	if err == nil {
		httpCfg.CSRFType = strValue
	}
	strValue, err = global.String("CSRFSource")
	if err == nil {
		httpCfg.CSRFSource = strValue
	}
	intValue, err = global.Int("CSRFExpire")
	if err == nil {
		httpCfg.CSRFExpire = intValue
	}
	strValue, err = global.String("Static")
	if err == nil {
		httpCfg.Static = strings.Split(strValue, ";")
		//路径规范化
		//统一规范为相对路径,相对于{Config.Root}所指定的目录
		// \content\css ==> content/css
		for i, p := range httpCfg.Static {
			httpCfg.Static[i] = NormalizePath(p)
		}
	}

	strValue, err = global.String("View")
	if err == nil {
		httpCfg.View = strValue
	}
	boolValue, err = global.Bool("Precompile")
	if err == nil {
		httpCfg.Precompile = boolValue
	}
	strValue, err = global.String("Api")
	if err == nil {
		httpCfg.Api = strValue
	}
	strValue, err = global.String("Favicon")
	if err == nil {
		httpCfg.Favicon = strValue
	}
	strValue, err = global.String("Robots")
	if err == nil {
		httpCfg.Robots = strValue
	}
	boolValue, err = global.Bool("Log")
	if err == nil {
		httpCfg.Log = boolValue
	}
	strValue, err = global.String("LogType")
	if err == nil {
		httpCfg.LogType = strValue
	}
	strValue, err = global.String("LogPath")
	if err == nil {
		httpCfg.LogPath = strValue
	}
	boolValue, err = global.Bool("LogAsync")
	if err == nil {
		httpCfg.LogAsync = boolValue
	}
	strValue, err = global.String("LayoutConfigPath")
	if err == nil {
		httpCfg.LayoutConfigPath = strValue
	}
	strValue, err = global.String("TemplateExt")
	if err == nil {
		httpCfg.TemplateExt = strValue
	}
	strValue, err = global.String("TemplateName")
	if err == nil {
		httpCfg.TemplateName = strValue
	}
	strValue, err = global.String("SessionCookieName")
	if err == nil {
		httpCfg.SessionCookieName = strValue
	}
	strValue, err = global.String("CSRFCookieName")
	if err == nil {
		httpCfg.CSRFCookieName = strValue
	}
	strValue, err = global.String("CSRFTokenName")
	if err == nil {
		httpCfg.CSRFTokenName = strValue
	}
	intValue, err = global.Int("MaxRequestMemory")
	if err == nil {
		httpCfg.MaxRequestMemory = intValue
	}

	//读取视图配置
	if httpCfg.LayoutConfigPath != "" {
		var content, err = ioutil.ReadFile(filepath.Join(appDir, httpCfg.LayoutConfigPath))
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal(content, &httpCfg.ViewConfig)
		if err != nil {
			return nil, err
		}
		//路径规范化,相对于{Config.Root}/{Config.View}所指定的目录
		for k, v := range httpCfg.ViewConfig.LayoutMap {
			httpCfg.ViewConfig.LayoutMap[k] = NormalizePath(v)
		}
		for k, v := range httpCfg.ViewConfig.LayoutSpec {
			delete(httpCfg.ViewConfig.LayoutSpec, k)
			httpCfg.ViewConfig.LayoutSpec[NormalizePath(k)] = v
		}
	}

	return httpCfg, nil
}

// NormalizePath 规范化路径
// 规范化后统一使用/分隔符,所有绝对路径全部变为相对路径
//
// 例如:
//  \test\test.html ==> test/test.html
func NormalizePath(path string) string {
	return strings.Trim(strings.Replace(path, "\\", "/", -1), "/")
}
