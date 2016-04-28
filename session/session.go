// Package session 实现了基于内存的session管理功能
package session

import "sync"

// Session
type Session interface {
	// SessionId 返回Session的Id
	SessionId() string
	// Value 获取值
	Value(key string) (interface{}, bool)
	// String 获取字符串
	String(key string) (string, bool)
	//  Int 获取整数值
	Int(key string) (int, bool)
	// Bool 获取bool值
	Bool(key string) (bool, bool)
	// Float 获取浮点值
	Float(key string) (float64, bool)
	// SetValue 设置值
	SetValue(key string, value interface{})
	// SetString 设置字符串
	SetString(key string, value string)
	// SetInt 设置整数值
	SetInt(key string, value int)
	// SetBool 设置bool值
	SetBool(key string, value bool)
	// SetFloat 设置浮点值
	SetFloat(key string, value float64)
	// Delete 删除指定键
	Delete(key string)
	// SetDeadline 设置有效期限
	SetDeadline(second int64)
	// Dead 让Session立即过期
	Die()
	// Dead 判断Session是否过期
	Dead() bool
}

// Session容器
type SessionContainer interface {
	// CreateSession 创建Session
	CreateSession() (Session, bool)
	// Session 获取Session
	Session(sessionId string) (Session, bool)
	// Close 关闭SessionProvider,关闭之后将无法使用
	Close()
	// Closed 确认当前SessionProvider是否已经关闭
	Closed() bool
}

// SessionContainer创建器
//  expire:session有效期
//  source:存储源
type SessionContainerCreator func(expire int64, source string) (SessionContainer, error)

var (
	mu       sync.Mutex                                 //互斥锁
	creators = make(map[string]SessionContainerCreator) //SessionContainer创建器映射
)

// NewSessionContainer 创建一个新的Session容器
//  expire:最大过期时间(秒)
func NewSessionContainer(kind string, expire int64, source string) (SessionContainer, error) {
	var creator, ok = creators[kind]
	if !ok {
		return nil, ErrorInvalidSessionKind.Format(kind).Error()
	}
	return creator(expire, source)
}

// Register 注册SessionContainer创建器
func Register(kind string, creator SessionContainerCreator) error {
	if creator == nil {
		panic(ErrorInvalidSessionProvider)
	}
	mu.Lock()
	defer mu.Unlock()
	creators[kind] = creator
	return nil
}
