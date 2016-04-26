package router

import (
	"reflect"
	"regexp"
	"strings"
)

// 基础路由
type BaseRouter struct {
	UnlimitedRouter
	namedChildren   map[string]Router //名称子路由
	unnamedChildren map[string]Router //非名称子路由
	reg             bool              //是否是正则路由
	regexp          *regexp.Regexp    //正则表达式
	keys            []string          //可提取keys
}

// NewBaseRouter 创建基本路由
//  name:路由名称,如果match为nil或者不包含正则,则使用name进行路由匹配
//  match:必须是指定规则的正则字符串或nil,格式可以为p{id=\d+}.html 解析为 ^p(\d+).html$
func NewBaseRouter(name string, match interface{}) (Router, error) {
	var r = new(BaseRouter)
	r.name = name
	if match != nil {
		var matchString, ok = match.(string)
		if !ok {
			return nil, ErrorInvalidMatchParam.Format(reflect.TypeOf(match).String(), "string").Error()
		}
		var seg, err = ParseReg(matchString)
		if err == nil {
			r.reg = true
			r.regexp = seg.Regexp
			r.keys = seg.Keys
		} else {
			r.reg = false
		}
	}
	r.namedChildren = make(map[string]Router, 0)
	r.unnamedChildren = make(map[string]Router, 0)
	r.preFilters = make([]PreFilter, 0)
	r.postFilters = make([]PostFilter, 0)
	return r, nil
}

// Named 返回当前是否使用Name进行路由匹配
func (this *BaseRouter) Named() bool {
	return !this.reg
}

// unifyName 统一名称
func (this *BaseRouter) unifyName(name string) string {
	return strings.ToLower(name)
}

// AddChild 添加子路由
func (this *BaseRouter) AddChild(router Router) {
	var name = this.unifyName(router.Name())
	var child, ok = this.Child(name)
	if ok {
		//合并路由
		child.AddChildren(router.Children())
	} else {
		//添加路由
		if router.Named() {
			this.namedChildren[name] = router
		} else {
			this.unnamedChildren[name] = router
		}
		router.SetParent(this)
	}
}

// AddChildren 批量添加子路由
func (this *BaseRouter) AddChildren(routers []Router) {
	for _, v := range routers {
		this.AddChild(v)
	}
}

// Child 返回指定名称的子路由
func (this *BaseRouter) Child(name string) (Router, bool) {
	name = this.unifyName(name)
	var r, ok = this.namedChildren[name]
	if ok {
		return r, ok
	}
	r, ok = this.unnamedChildren[name]
	return r, ok
}

// Children 返回全部子路由
func (this *BaseRouter) Children() []Router {
	var routers = make([]Router, len(this.namedChildren)+len(this.unnamedChildren))
	for _, v := range this.namedChildren {
		routers = append(routers, v)
	}
	for _, v := range this.unnamedChildren {
		routers = append(routers, v)
	}
	return []Router{}
}

// RemoveChild 移除指定名称的路由,并返回该路由
func (this *BaseRouter) RemoveChild(name string) (Router, bool) {
	name = this.unifyName(name)
	var r, ok = this.namedChildren[name]
	if ok {
		delete(this.namedChildren, name)
		return r, ok
	}
	r, ok = this.unnamedChildren[name]
	if ok {
		delete(this.unnamedChildren, name)
		return r, ok
	}
	return nil, false
}

// Match 匹配指定路由上下文,匹配成功则返回RouterExcutor
func (this *BaseRouter) Match(context RouterContext) (RouterExcutor, bool) {
	var segs = context.Segments()
	var data [][]string
	if !this.Named() {
		data = this.regexp.FindAllStringSubmatch(segs[0], 1)
		if len(data) <= 0 {
			return nil, false
		}
	}
	context.Match(1)
	defer context.Unmatch(1)

	segs = context.Segments()
	var executor RouterExcutor
	if len(segs) <= 0 {
		if this.executorGenerator == nil {
			//匹配失败
			return nil, false
		}
		//成功匹配
		executor = this.executorGenerator()
		executor.SetRouter(this)
		executor.SetRouterContext(context)
	} else {
		//路由传递
		var name = this.unifyName(segs[0])
		var r, ok = this.namedChildren[name]
		if ok {
			executor, ok = r.Match(context)
		}
		if !ok {
			for _, v := range this.unnamedChildren {
				executor, ok = v.Match(context)
				if ok {
					break
				}
			}
		}
		if !ok {
			return nil, false
		}
	}
	//匹配成功设置路由值
	if len(data) >= 1 {
		var values = data[0][1:]
		for i, v := range values {
			context.SetValue(this.keys[i], v)
		}
	}
	return executor, true
}
