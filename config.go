package tinygo

var LayoutConfigPath = "views/layout.json" // 布局配置文件名
var TemplateExt = "html"                   // 视图文件扩展名
var TemplateName = "Content"               // 模板文件内模板名,用于返回部分视图时使用
var SessionCookieName = "ssid"             // Session Cookie名
var CSRFCookieName = "csrfid"              // CSRF Cookie 名
var CSRFTokenName = "csrf"                 // CSRF 表单名
var MaxMemory = 32 << 20                   // 最大文件表单占用内存大小32 MB
var App string = "Default"                 // 应用名称
var Mode string = "debug"                  // 启动模式,可以为debug或release
var Https bool = false                     // 是否启用https,可选,默认为false
var Port uint16 = 80                       // 监听端口,可选,默认为80，https为true则默认为443
var Cert string                            // 证书(PEM)路径,如果启用了https则必填
var PrivateKey string                      // 私钥(PEM)路径,如果启用了https则必填
var Home string                            // 首页地址
var Session bool = true                    // 是否启用session
var Sessiontype string = "memory"          // session类型,参考tinygo/session,默认为memory
var SessionExpire int64 = 1800             // session过期时间,单位为秒
var Csrf bool = false                      // 是否启用csrf
var Csrfexpire int64                       // csrf token过期时间
var Static []string = []string{"content"}  // 静态文件目录,默认为"content"
var View string = "views"                  // 视图文件目录,默认为"views"
var Precompile bool = false                // 是否预编译页面路径,默认为false
var Api string = "json"                    // 使用Api返回的数据的解析格式,默认为auto(其他设置包括json,xml)

// readConfig 读取配置文件
func readConfig() error {
	return nil
}
