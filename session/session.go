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

// Session提供器
type SessionProvider interface {
	// CreateSession 创建Session
	CreateSession() (Session, bool)
	// Session 获取Session
	Session(sessionId string) (Session, bool)
	// Clean 清理过期Session
	Clean()
}

// SessionProvider创建器
type SessionProviderCreator func(expire int64) (SessionProvider, error)

var (
	mu       sync.Mutex                                     //生成器互斥锁
	creators = make(map[SessionType]SessionProviderCreator) //配置解析器
)

// NewConfig 创建一个新的Config
//  path:配置文件路径
func NewSessionProvider(kind SessionType, expire int64) (SessionProvider, error) {
	var creator, ok = creators[kind]
	if !ok {
		return nil, SessionErrorInvalidSessionKind.Format(kind).Error()
	}
	return creator(expire)
}

// RegisterSessionProviderCreator 注册SessionProvider创建器
func RegisterSessionProviderCreator(kind SessionType, creator SessionProviderCreator) error {
	mu.Lock()
	defer mu.Unlock()
	creators[kind] = creator
	return nil
}
