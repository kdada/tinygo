package connector

import (
	"net/http"
	"strings"
)

// Http上下文
type HttpContext struct {
	Segments       []string            //路由段
	Request        *http.Request       //http请求
	ResponseWriter http.ResponseWriter //http响应
}

// Http处理器
type HttpHandler struct {
	dispatcher Dispatcher
}

// ServeHTTP 处理http请求
func (this *HttpHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	var context = &HttpContext{nil, r, rw}
	context.Segments = strings.Split(r.URL.Path, `/`)
	var i = len(context.Segments) - 1
	for ; i > 0; i-- {
		if context.Segments[i] != "" {
			break
		}
	}
	context.Segments = context.Segments[:i+1]
	this.dispatcher.Dispatch(context.Segments, context)
}

// Http连接器
type HttpConnector struct {
	server     *http.Server //http服务
	addr       string       //监听地址
	dispatcher Dispatcher
}

// NewHttpConnector 创建Http连接器
func NewHttpConnector(source string) (Connector, error) {
	return &HttpConnector{nil, source, nil}, nil
}

// Init 初始化连接器设置
func (this *HttpConnector) Init() error {
	return nil
}

// Run 运行(接受连接并进行处理,阻塞)
func (this *HttpConnector) Run() error {
	if this.dispatcher == nil {
		panic(ErrorInvalidDispatcher.Format("http"))
	}
	this.server = &http.Server{Addr: this.addr, Handler: &HttpHandler{this.dispatcher}}
	return this.server.ListenAndServe()
}

// Stop 停止运行
func (this *HttpConnector) Stop() error {
	return ErrorFailToStop.Format("http").Error()
}

// Dispatcher 返回当前调度器
func (this *HttpConnector) Dispatcher() Dispatcher {
	return this.dispatcher
}

// SetDispatcher 设置调度器
func (this *HttpConnector) SetDispatcher(dispatcher Dispatcher) {
	this.dispatcher = dispatcher
}
