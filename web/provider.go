package web

import (
	"reflect"

	"github.com/kdada/tinygo/meta"
)

// http上下文值提供器
type ContextValueProvider struct {
	Context *Context
	Finder  ContextValueFinder
	Name    string
	Type    reflect.Type
}

// String 根据名称和类型返回相应的字符串值
func (this *ContextValueProvider) String() []string {
	return this.Finder.String(this.Context, this.Name, this.Type)
}

// Value 根据名称和类型生成相应类型的数据
func (this *ContextValueProvider) Value() interface{} {
	return this.Finder.Value(this.Context, this.Name, this.Type)
}

// http上下文值容器
// 优先级:web.Processor.Finders > web.Processor.MutiTypeFinders > meta.GlobalValueContainer
type ContextValueContainer struct {
	Context *Context
}

// NewContextValueContainer 创建http上下文值容器
func NewContextValueContainer(context *Context) *ContextValueContainer {
	return &ContextValueContainer{
		context,
	}
}

// String 根据名称和类型返回相应的字符串值,返回的bool表示该值是否存在
func (this *ContextValueContainer) Contains(name string, t reflect.Type) (meta.ValueProvider, bool) {
	// 查找 web.Processor.Finders
	var finder, ok = this.Context.Processor.Finders[t.String()]
	if ok {
		ok = finder.Contains(this.Context, name, t)
	}
	if !ok {
		// 查找 web.Processor.MutiTypeFinders
		for _, f := range this.Context.Processor.MutiTypeFinders {
			if f.Contains(this.Context, name, t) {
				finder = f
				ok = true
				break
			}
		}
	}
	if ok {
		return &ContextValueProvider{
			this.Context,
			finder,
			name,
			t,
		}, true
	}
	return meta.GlobalValueContainer.Contains(name, t)
}
