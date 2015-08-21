package session

// Session
type ISession interface {
	// SessionId 返回Session的Id
	SessionId() string
	// Value 获取值
	Value(key string) (interface{}, bool)
	// String 获取字符串
	String(key string) (string, bool)
	//  Int 获取整数值
	Int(key string) (int64, bool)
	// Bool 获取bool值
	Bool(key string) (bool, bool)
	// Float 获取浮点值
	Float(key string) (float64, bool)
	// SetValue 设置值
	SetValue(key string, value interface{})
	// SetString 设置字符串
	SetString(key string, value string)
	// SetInt 设置整数值
	SetInt(key string, value int64)
	// SetBool 设置bool值
	SetBool(key string, value bool)
	// SetFloat 设置浮点值
	SetFloat(key string, value float64)
	// SetDeadline 设置有效期限
	SetDeadline(second int64)
	// Dead 让Session立即过期
	Die()
	// Dead 判断Session是否过期
	Dead() bool
}

// Session提供器
type ISessionProvider interface {
	// CreateSession 创建Session
	CreateSession() (ISession, bool)
	// Session 获取Session
	Session(sessionId string) (ISession, bool)
	// Clean 清理过期Session
	Clean()
}
