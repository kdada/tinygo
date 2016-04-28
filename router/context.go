package router

type BaseContext struct {
	Segs  []string // 路由段信息
	Level int      //当前路由级别
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

// Value 返回路由值
func (this *BaseContext) Value(name string) (string, bool) {
	return "", false
}

// SetValue 设置路由值
func (this *BaseContext) SetValue(name string, value string) {

}

// Data 返回路由上下文携带的信息
func (this *BaseContext) Data() interface{} {
	return nil
}
