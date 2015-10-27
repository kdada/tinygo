#TinyGo
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


一个tinygo项目结构如下:
```
main.go
web.cfg
|-views
  |-layout.json  
  |-layout
    |-layout.html
  |-home
    |-index.html
|-content
  |-css
  |-js
  |-img
|-controllers
  |-home.go
|-models
  |-model.go
|-routers
  |-routers.go
|-services
  |-service.go
```
(1) main.go为项目启动文件  
(2) views内放置视图(页面)文件,目录名应该与控制器名(去掉Controller的部分)相同.  
当然也可以与控制器不同,但是此时就需要在控制器内指定视图文件的路径.  
layout.json为视图布局配置文件  
(3) content为静态文件目录,可以通过/content/js等路径访问  
(4) controllers为控制器文件目录,文件名为(去掉Controller的部分),一个控制器一个文件  
(5) models为模型(Model)目录  
(6) routers目录[非必须]为路由目录,控制器路由可以在routers.go中进行注册  
(7) services目录[非必须]为服务目录,网络服务方法,数据库服务方法等可以放在该目录中



web.cfg 配置文件范例如下
```ini
#tinygo配置文件

#启动模式,可以为debug或release
mode = debug

#是否启用https,可选,默认为false
https = false

#监听端口,可选,默认为80，https为true则默认为443
port = 80

#证书(PEM)路径,如果启用了https则必填
#cert = keys/cert.pem

#私钥(PEM)路径,如果启用了https则必填
#pkey = keys/privatekey.pem

#首页
home = /home/index

#是否启用session
session = true

#session类型,参考tinygo/session,默认为memory
sessiontype = memory

#session过期时间,单位为秒
sessionexpire = 600

#静态文件目录,默认为"content",多个目录用;分隔
static = content

#视图文件目录,默认为"views"
view = views

#默认错误页面路径,默认为空,该路径需要在view所指定的视图文件目录中
pageerr = errors/error.html

#是否预编译页面路径,默认为false,发布模式下最好为true以提高效率
precompile = false

#使用Api返回的数据的解析格式,默认为auto
api = json
```
在该配置文件中可以更改相应的目录名称和其他http设置  
在产品环境下,mode应该设置为release,precompile应该设置为true

layout.json 布局配置文件范例如下
```json
{
	"LayoutMap":{
		"Default":"layout/layout.html"
		"Index":"layout/layout.html"
	},
	"DefaultLayout":"Default",
	"LayoutSpec":{
		"home/index.html":"Index"
		"home/login.html":"Default"
		"xxhome/":"Index"
	}
}
```
(1) LayoutMap定义了布局文件和布局文件别名,按 别名:布局文件路径 进行映射  
(2) LayoutSpec定义了视图文件或目录下的所有视图文件使用相应的布局文件  
直接定义的视图文件到布局的映射优先级高于目录到布局的映射  
(3) DefaultLayout定义了LayoutSpec没有定义的视图文件使用的默认布局  
  
  
  
  
  
开发中:  
(1) 过期Session清理(是否需要暂停所有Session获取?)  
(2) CSRF  
(3) template func  
