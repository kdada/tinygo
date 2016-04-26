package connector

import "sync"

// 连接器
type Connector interface {
	// Init 初始化连接器设置
	Init() error
	// Run 运行(接受连接并进行处理,阻塞)
	Run()
	// Stop 停止运行
	Stop() error
	// EventHandler 返回事件处理器
	EventHandler() EventHandler
	// SetHandler 设置事件处理器
	SetEventHandler(handler EventHandler)
}

// 事件处理器
type EventHandler interface {
	// 连接器收到新连接后触发
	Connected(context interface{})
	// 出现错误时触发
	Error(code int, context interface{})
}

// 链接器创建器
//  param: 日志参数
type ConnectorCreator func(param interface{}) (Connector, error)

var (
	mu       sync.Mutex                          //互斥锁
	creators = make(map[string]ConnectorCreator) //日志创建器映射
)

// NewConnector 创建一个新的Connector
//  kind:日志类型
func NewConnector(kind string, param interface{}) (Connector, error) {
	var creator, ok = creators[kind]
	if !ok {
		return nil, ErrorInvalidKind.Format(kind).Error()
	}
	return creator(param)
}

// Register 注册Connector创建器
func Register(kind string, creator ConnectorCreator) {
	if creator == nil {
		panic(ErrorInvalidConnectorCreator)
	}
	mu.Lock()
	defer mu.Unlock()
	creators[kind] = creator
}
