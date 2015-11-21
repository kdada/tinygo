package session

import (
	"sync"
	"time"
)

// 内存Session,每次操作都会使Session有效时间延长
type MemSession struct {
	provider  *MemSessionProvider    //会话提供器
	sessionId string                 //会话id
	data      map[string]interface{} //数据
	deadline  int64                  //死亡时间(秒),从1970年开始
}

// newMemSession 创建内存Session
func newMemSession(sessionId string) *MemSession {
	var ss = new(MemSession)
	ss.sessionId = sessionId
	ss.data = make(map[string]interface{}, 0)
	return ss
}

// SessionId 返回Session的Id
func (this *MemSession) SessionId() string {
	return this.sessionId
}

// Value 获取值
func (this *MemSession) Value(key string) (interface{}, bool) {
	v, ok := this.data[key]
	this.SetDeadline(this.provider.defaultExpire)
	return v, ok
}

// String 获取字符串
func (this *MemSession) String(key string) (string, bool) {
	v, ok := this.data[key]
	s, ok := v.(string)
	this.SetDeadline(this.provider.defaultExpire)
	return s, ok
}

//  Int 获取整数值
func (this *MemSession) Int(key string) (int, bool) {
	v, ok := this.data[key]
	s, ok := v.(int)
	this.SetDeadline(this.provider.defaultExpire)
	return s, ok
}

// Bool 获取bool值
func (this *MemSession) Bool(key string) (bool, bool) {
	v, ok := this.data[key]
	s, ok := v.(bool)
	this.SetDeadline(this.provider.defaultExpire)
	return s, ok
}

// Float 获取浮点值
func (this *MemSession) Float(key string) (float64, bool) {
	v, ok := this.data[key]
	s, ok := v.(float64)
	this.SetDeadline(this.provider.defaultExpire)
	return s, ok
}

// SetValue 设置值
func (this *MemSession) SetValue(key string, value interface{}) {
	this.data[key] = value
	this.SetDeadline(this.provider.defaultExpire)
}

// SetString 设置字符串
func (this *MemSession) SetString(key string, value string) {
	this.data[key] = value
	this.SetDeadline(this.provider.defaultExpire)
}

// SetInt 设置整数值
func (this *MemSession) SetInt(key string, value int) {
	this.data[key] = value
	this.SetDeadline(this.provider.defaultExpire)
}

// SetBool 设置bool值
func (this *MemSession) SetBool(key string, value bool) {
	this.data[key] = value
	this.SetDeadline(this.provider.defaultExpire)
}

// SetFloat 设置浮点值
func (this *MemSession) SetFloat(key string, value float64) {
	this.data[key] = value
	this.SetDeadline(this.provider.defaultExpire)
}

// Delete 删除键
func (this *MemSession) Delete(key string) {
	delete(this.data, key)
	this.SetDeadline(this.provider.defaultExpire)
}

// SetDeadline 设置有效期限
//  second:从当前时间开始有效的秒数
func (this *MemSession) SetDeadline(second int64) {
	this.deadline = time.Now().Unix() + second
}

// Die 让Session立即过期
func (this *MemSession) Die() {
	this.deadline = 0
}

// Dead 判断Session是否过期
func (this *MemSession) Dead() bool {
	return time.Now().Unix() > this.deadline
}

// 内存Session提供器
type MemSessionProvider struct {
	sessionCounter int                    //session计数器
	sessions       map[string]*MemSession //存储Session
	defaultExpire  int64                  //默认过期时间
	rwm            sync.RWMutex           //读写锁
}

// newMemSessionProvider 创建Session提供器
func newMemSessionProvider(expire int64) (SessionProvider, error) {
	var provider = new(MemSessionProvider)
	provider.sessions = make(map[string]*MemSession, 100)
	provider.defaultExpire = expire
	return provider, nil
}

// CreateSession 创建Session
func (this *MemSessionProvider) CreateSession() (Session, bool) {
	this.rwm.Lock()
	defer this.rwm.Unlock()
	var sessionId = Guid()
	var ss = newMemSession(sessionId)
	ss.provider = this
	ss.SetDeadline(this.defaultExpire)
	this.sessionCounter++
	this.sessions[sessionId] = ss
	return ss, true
}

// Session 获取Session
func (this *MemSessionProvider) Session(sessionId string) (Session, bool) {
	this.rwm.RLock()
	defer this.rwm.RUnlock()
	var ss, ok = this.sessions[sessionId]
	if ok && ss.Dead() {
		delete(this.sessions, sessionId)
		return nil, false
	}
	return ss, ok
}

// Clean 清理过期Session
func (this *MemSessionProvider) Clean() {
	this.rwm.Lock()
	defer this.rwm.Unlock()
	for k, v := range this.sessions {
		if v.Dead() {
			delete(this.sessions, k)
		}
	}
}
