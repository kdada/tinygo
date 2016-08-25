package session

import (
	"sync"
	"time"

	"github.com/kdada/tinygo/util"
)

// 内存Session,每次操作都会使Session有效时间延长
type MemSession struct {
	sessionId string                 //会话id
	data      map[string]interface{} //数据
	deadline  int                    //死亡时间(秒),从1970年开始
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
	return v, ok
}

// String 获取字符串
func (this *MemSession) String(key string) (string, bool) {
	v, ok := this.data[key]
	s, ok := v.(string)
	return s, ok
}

//  Int 获取整数值
func (this *MemSession) Int(key string) (int, bool) {
	v, ok := this.data[key]
	s, ok := v.(int)
	return s, ok
}

// Bool 获取bool值
func (this *MemSession) Bool(key string) (bool, bool) {
	v, ok := this.data[key]
	s, ok := v.(bool)
	return s, ok
}

// Float 获取浮点值
func (this *MemSession) Float(key string) (float64, bool) {
	v, ok := this.data[key]
	s, ok := v.(float64)
	return s, ok
}

// SetValue 设置值
func (this *MemSession) SetValue(key string, value interface{}) {
	this.data[key] = value
}

// SetString 设置字符串
func (this *MemSession) SetString(key string, value string) {
	this.data[key] = value
}

// SetInt 设置整数值
func (this *MemSession) SetInt(key string, value int) {
	this.data[key] = value
}

// SetBool 设置bool值
func (this *MemSession) SetBool(key string, value bool) {
	this.data[key] = value
}

// SetFloat 设置浮点值
func (this *MemSession) SetFloat(key string, value float64) {
	this.data[key] = value
}

// Delete 删除键
func (this *MemSession) Delete(key string) {
	delete(this.data, key)
}

// SetDeadline 设置有效期限
//  second:从当前时间开始有效的秒数
func (this *MemSession) SetDeadline(second int) {
	this.deadline = int(time.Now().Unix()) + second
}

// Deadline 获取Session过期时间
func (this *MemSession) Deadline() int {
	return this.deadline
}

// Die 让Session立即过期
func (this *MemSession) Die() {
	this.deadline = 0
}

// Dead 判断Session是否过期
func (this *MemSession) Dead() bool {
	return int(time.Now().Unix()) > this.deadline
}

// 内存Session容器
type MemSessionContainer struct {
	sessionCounter int                    //session计数器
	sessions       map[string]*MemSession //存储Session
	defaultExpire  int                    //默认过期时间
	rwm            sync.RWMutex           //读写锁
	closed         bool                   //是否关闭
}

// NewMemSessionContainer 创建Session提供器(数据存储在内存中,source参数无效)
func NewMemSessionContainer(expire int, source string) (SessionContainer, error) {
	var container = new(MemSessionContainer)
	container.sessions = make(map[string]*MemSession, 100)
	container.defaultExpire = expire
	container.closed = false
	go func() {
		//Clean [60,expire/2]
		//每分钟检查一次,达到指定时间时清理一次Session并计算下一次清理的时间
		var minCleanSep = time.Duration(60) * time.Second
		var maxCleanSep = time.Duration(expire/2) * time.Second
		if maxCleanSep < minCleanSep {
			maxCleanSep = minCleanSep
		}
		var rangeCleanSep = float32(maxCleanSep - minCleanSep)
		var cleanSep = minCleanSep
		var passTime = time.Duration(0)
		for !container.closed {
			time.Sleep(minCleanSep)
			passTime += minCleanSep
			if !container.closed && passTime >= cleanSep {
				var dead = container.Clean()
				//计算下一次Clean的时间间隔
				passTime = 0
				var alive = len(container.sessions)
				var total = float32(dead + alive)
				if total <= 0 {
					total = 1.0
				}
				var deadRate = float32(dead) / total
				if deadRate >= 0.2 {
					//本次清理超过20%
					cleanSep = minCleanSep
				} else {
					deadRate *= 5
					cleanSep = minCleanSep + time.Duration(rangeCleanSep*deadRate)
				}
			}
		}
	}()
	return container, nil
}

// CreateSession 创建Session
func (this *MemSessionContainer) CreateSession() (Session, bool) {
	if this.closed {
		return nil, false
	}
	this.rwm.Lock()
	defer this.rwm.Unlock()
	var sessionId = util.NewUUID().Hex()
	var ss = newMemSession(sessionId)
	ss.SetDeadline(this.defaultExpire)
	this.sessionCounter++
	this.sessions[sessionId] = ss
	return ss, true
}

// Session 获取Session
func (this *MemSessionContainer) Session(sessionId string) (Session, bool) {
	if this.closed {
		return nil, false
	}
	this.rwm.RLock()
	defer this.rwm.RUnlock()
	var ss, ok = this.sessions[sessionId]
	if ok {
		if ss.Dead() {
			delete(this.sessions, sessionId)
			return nil, false
		}
		//更新Session的过期时间
		ss.SetDeadline(this.defaultExpire)
	}
	return ss, ok
}

// Clean 清理过期Session 并返回清理的数量
func (this *MemSessionContainer) Clean() int {
	if this.closed {
		return 0
	}
	this.rwm.Lock()
	defer this.rwm.Unlock()
	var count = 0
	for k, v := range this.sessions {
		if v.Dead() {
			count++
			delete(this.sessions, k)
		}
	}
	return count
}

// Close 关闭SessionProvider,关闭之后将无法使用
func (this *MemSessionContainer) Close() {
	this.closed = true
}

// Closed 确认当前SessionProvider是否已经关闭
func (this *MemSessionContainer) Closed() bool {
	return this.closed
}
