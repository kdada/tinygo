package router

import "regexp"

// 基础路由上下文
type BaseContext struct {
	Segs  []string //路由段信息
	Level int      //当前路由级别
}

// 分隔符正则表达式
var spReg = regexp.MustCompile(`[\\/]+`)

// NewBaseContext 使用path创建基础路由上下文,path可以用\或/分割
func NewBaseContext(path string) *BaseContext {
	var segs = spReg.Split(path+"/", -1)
	segs = segs[:len(segs)-1]
	return &BaseContext{
		segs,
		0,
	}
}

// Segments 返回可匹配路由段
func (this *BaseContext) Segments() []string {
	if this.Level > len(this.Segs) {
		return []string{}
	}
	return this.Segs[this.Level:]
}

// AllSegments 返回所有路由段
func (this *BaseContext) AllSegments() []string {
	return this.Segs
}

// Match 匹配数量
func (this *BaseContext) Match(count int) {
	this.Level += count
}

// Unmatch 失配数量
func (this *BaseContext) Unmatch(count int) {
	this.Level -= count
}

// Matched 返回当前已匹配的路由段数量(即经过的路由数量)
func (this *BaseContext) Matched() int {
	return this.Level
}

// Pure 返回当前是否未匹配任何路由
func (this *BaseContext) Pure() bool {
	return this.Level == 0
}

// Value 返回路由值
func (this *BaseContext) Value(name string) (string, bool) {
	return "", false
}

// SetValue 设置路由值
func (this *BaseContext) SetValue(name string, value string) {

}
