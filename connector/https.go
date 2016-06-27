package connector

import (
	"net/http"
	"strings"
)

// Https连接器
type HttpsConnector struct {
	HttpConnector
	cert string
	key  string
}

// NewHttpsConnector 创建Http连接器
//  source:格式如: 127.0.0.1:8080;Cert=path/to/cert;Key=path/to/key
func NewHttpsConnector(source string) (Connector, error) {
	var sources = strings.Split(source, ";")
	if len(sources) <= 0 {
		return nil, ErrorParamNotFound.Format("ip地址").Error()
	}
	var addr = sources[0]
	var info = make(map[string]string, 2)
	for i := 1; i < len(sources); i++ {
		var kv = strings.Split(sources[i], "=")
		if len(kv) == 2 {
			info[strings.ToLower(strings.TrimSpace(kv[0]))] = strings.TrimSpace(kv[1])
		}
	}
	var c = new(HttpsConnector)
	c.addr = addr
	var ok bool = false
	c.cert, ok = info["cert"]
	if !ok {
		return nil, ErrorParamNotFound.Format("Cert").Error()
	}
	c.key = info["key"]
	if !ok {
		return nil, ErrorParamNotFound.Format("Key").Error()
	}
	return c, nil
}

// Run 运行(接受连接并进行处理,阻塞)
func (this *HttpsConnector) Run() error {
	if this.dispatcher == nil {
		panic(ErrorInvalidDispatcher.Format("https"))
	}
	this.server = &http.Server{Addr: this.addr, Handler: &HttpHandler{this.dispatcher}}
	return this.server.ListenAndServeTLS(this.cert, this.key)
}

// Stop 停止运行
func (this *HttpsConnector) Stop() error {
	return ErrorFailToStop.Format("https").Error()
}
