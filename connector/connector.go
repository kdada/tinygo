//  Package connector 实现了基本的连接器接口
package connector

import "sync"

// Dispatcher 调度器
type Dispatcher interface {
	// Dispatch 分发
	//  segments:用于进行分发的路径段信息
	//  data:连接携带的数据
	Dispatch(segments []string, data interface{})
}

// 连接器
type Connector interface {
	// Init 初始化连接器设置
	Init() error
	// Run 运行(接受连接并进行处理,阻塞)
	Run() error
	// Stop 停止运行
	Stop() error
	// Dispatcher 返回当前调度器
	Dispatcher() Dispatcher
	// SetDispatcher 设置调度器
	SetDispatcher(dispatcher Dispatcher)
}

// 连接器创建器
//  suorce: 连接器监听位置(例如:127.0.0.1:9999表示监听127.0.0.1上的9999端口)
type ConnectorCreator func(source string) (Connector, error)

var (
	mu       sync.Mutex                          //互斥锁
	creators = make(map[string]ConnectorCreator) //日志创建器映射
)

// NewConnector 创建一个新的Connector
//  kind:日志类型
func NewConnector(kind string, source string) (Connector, error) {
	var creator, ok = creators[kind]
	if !ok {
		return nil, ErrorInvalidKind.Format(kind).Error()
	}
	return creator(source)
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
