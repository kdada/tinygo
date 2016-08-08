#TinyGo 1.0
一个轻量级Http框架  

查看tinygo文档的方法如下:  
(1) 确保GOROOT和GOPATH环境变量设置正确  
(2) 确保tinygo的位置在$GOPATH/src/github.com/kdada/tinygo  
(3) 在shell中执行如下命令  
```bash
cd $GOROOT/bin
godoc -http :9999
```
(4) 通过 http://localhost:9999 访问本地文档,即可在Packages中找到tinygo的文档  


一个tinygo项目基本结构如下(可通过web.cfg变更项目结构):  
```
main.go
web.cfg
|-app
  |-favicon.ico
  |-robots.txt
  |-css
  |-js
  |-img
  |-views
    |-layout.json
    |-layout
      |-layout.html
    |-home
      |-index.html
|-controllers
  |-home.go
|-models
  |-model.go
|-services
  |-service.go
```
(1) main.go为项目启动文件  
(2) app目录为应用文件目录,用于存放前端所有文件  
(3) views内放置视图页面文件,layout.json为视图布局配置文件  
(4) controllers为控制器文件目录,一般一个控制器一个文件  
(5) models为模型文件目录  
(6) services为服务文件目录,网络服务方法,数据库服务方法等可以放在该目录中  



web.cfg 配置文件范例如下(UTF-8格式)  
```ini

#tinygo配置文件

#应用名称
App = blog

#是否启用https,可选,默认为false
Https = false

#监听端口,默认为80，https为true则默认为443
Port = 8080

#证书(PEM)路径,如果启用了https则必填
#Cert = keys/cert.pem

#私钥(PEM)路径,如果启用了https则必填
#PrivateKey = keys/privatekey.pem

#首页
Home = /home/index

#是否启用session
Session = true

#session类型,参考tinygo/session,默认为memory
SessionType = memory

#session源,参考tinygo/session,默认为空
#SessionSource =

#session过期时间,单位为秒,csrf的session过期时间也由该值决定
SessionExpire = 1800

#Session Cookie名,默认为ssid
SessionCookieName = blog_ssid

#Session Cookie的过期时间,单位为秒,默认为0(0表示浏览器关闭后过期)
#SessionCookieExpire

#是否启用CSRF,默认为false
CSRF = true

#CSRF session类型,参考tinygo/session,默认为memory
CSRFType = memory

#CSRF session源,参考tinygo/session,默认为空
#CSRFSource =

#CSRF的token过期时间,单位为秒
#该过期时间不是CSRF的session的过期时间,而是生成的每个CSRF token的过期时间
CSRFExpire = 300

#CSRF Cookie名,默认为xid
CSRFCookieName = blog_xid

#CSRF 表单字段名,默认为csrf
CSRFTokenName = __xfield

#CSRF Cookie的过期时间,单位为秒,默认为0(0表示浏览器关闭后过期)
#CSRFCookieExpire

#静态文件目录,默认为content,多个目录用;分隔,最后一级目录名不能重复
Static = app/js;app/css;app/tmpl

#静态文件目录是否允许显示目录列表,默认为false
List = true

#视图文件目录,默认为views
View = app/views

#布局配置文件名,默认为空
LayoutConfigPath = app/views/layout.json

#视图文件扩展名,默认为html
TemplateExt = html

#模板文件内部分模板名,用于返回部分视图时使用,默认为Content
TemplateName = Content

#是否预编译页面路径,默认为false,生产环境为true可以提高效率
Precompile = false

#使用Api返回的数据的解析格式,默认为auto(其他设置包括json,xml)
Api = json

#网站图标路径,默认为favicon.ico
Favicon = app/favicon.ico

#爬虫协议文件路径,默认为robots.txt
Robots = app/robots.txt

#是否启用日志,默认为true
Log = true

#日志类型,参考tinygo/log,默认为console
LogType = console

#日志路径,日志类型为file的时候需要设置
#LogPath = logs

#异步日志,默认为false
LogAsync = false

#单次请求最大占用内存大小,默认32 MB
MaxRequestMemory = 33554432

```

layout.json 布局配置文件范例如下(UTF-8格式)  
```json
{
	"LayoutMap":{
		"Default":"layout/layout.html",
		"Index":"layout/layout.html"
	},
	"DefaultLayout":"Default",
	"LayoutSpec":{
		"home/index.html":"Index",
		"home/login.html":"Default",
		"xxhome/":"Index"
	}
}
```

(1) LayoutMap定义了布局文件和布局文件别名,按 别名:布局文件路径 进行映射  
(2) LayoutSpec定义了视图文件或目录下的所有视图文件使用相应的布局文件,直接定义的视图文件到布局的映射优先级高于目录到布局的映射  
(3) DefaultLayout定义了LayoutSpec没有定义的视图文件使用的默认布局,仅对非布局文件有效  
(4) 使用该配置文件可以实现多重布局  

