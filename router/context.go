package router

type BaseContext struct {
	Segs    []string // 路由段信息
	Level   int      //当前路由级别
	Stopped bool     //是否已经停止路由
}

// Segments 返回可匹配路由段
func (this *BaseContext) Segments() []string {
	if this.Level > len(this.Segs) {
		return []string{}
	}
	return this.Segs[this.Level:]
}

// Match 匹配数量
func (this *BaseContext) Match(count int) {
	this.Level += count
}

// Unmatch 失配数量
func (this *BaseContext) Unmatch(count int) {
	this.Level -= count
}

// Terminated 路由过程是否终止
func (this *BaseContext) Terminated() bool {
	return this.Stopped
}

// Terminate 终止路由过程,终止后该路由上下文将不会被继续路由并且不会被执行器执行
func (this *BaseContext) Terminate() {
	this.Stopped = true
}

// Value 返回路由值
func (this *BaseContext) Value(name string) string {
	return ""
}

// SetValue 设置路由值
func (this *BaseContext) SetValue(name string, value string) {

}

// Data 返回路由上下文携带的信息
func (this *BaseContext) Data() interface{} {
	return nil
}
