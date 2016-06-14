package router

import (
	"reflect"
	"regexp"
	"strings"
)

// 基础路由
type BaseRouter struct {
	UnlimitedRouter
	children         map[string]Router
	normalChildren   map[string]Router //通常子路由
	abnormalChildren map[string]Router //非通常子路由
	match            string            //用于匹配的字符串,可能包含正则信息
	reg              bool              //是否是正则路由
	regexp           *regexp.Regexp    //正则表达式
	keys             []string          //可提取keys
}

// NewBaseRouter 创建基本路由,非正则路由不区分大小写,正则路由是否区分大小写由正则表达式确定
//  name:路由名称,如果match为nil,则使用name进行路由匹配
//  match:用于进行匹配的值,可以包含指定规则的正则字符串,格式可以为p{id=\d+}.html 解析为 ^p(\d+).html$
func NewBaseRouter(name string, match interface{}) (Router, error) {
	var r = new(BaseRouter)
	r.name = name
	if match != nil {
		var matchString, ok = match.(string)
		if !ok {
			return nil, ErrorInvalidMatchParam.Format(reflect.TypeOf(match).String(), "string").Error()
		}
		r.match = matchString
		var seg, err = ParseReg(matchString)
		if err == nil {
			r.reg = true
			r.regexp = seg.Regexp
			r.keys = seg.Keys
		} else {
			r.reg = false
		}
	} else {
		r.match = r.name
	}
	r.children = make(map[string]Router, 0)
	r.normalChildren = make(map[string]Router, 0)
	r.abnormalChildren = make(map[string]Router, 0)
	r.preFilters = make([]PreFilter, 0)
	r.postFilters = make([]PostFilter, 0)
	r.self = r
	return r, nil
}

// MatchString 返回当前路由用于进行匹配的字符串
func (this *BaseRouter) MatchString() string {
	return this.match
}

// Normal 返回当前路由是否为通常路由,通常路由可以使用MatchString()返回的字符串进行相等匹配
func (this *BaseRouter) Normal() bool {
	return !this.reg
}

// unify 统一字符串大小写
func (this *BaseRouter) unify(str string) string {
	return strings.ToLower(str)
}

// AddChild 添加子路由
func (this *BaseRouter) AddChild(router Router) {
	var child, ok = this.Child(router.Name())
	if ok {
		//合并路由
		child.AddChildren(router.Children())
	} else {
		//添加路由
		if router.Normal() {
			this.normalChildren[this.unify(router.MatchString())] = router
		} else {
			this.abnormalChildren[router.Name()] = router
		}
		this.children[router.Name()] = router
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
	var r, ok = this.children[name]
	return r, ok
}

// Children 返回全部子路由
func (this *BaseRouter) Children() []Router {
	var routers = make([]Router, len(this.children))
	for _, v := range this.children {
		routers = append(routers, v)
	}
	return []Router{}
}

// RemoveChild 移除指定名称的路由,并返回该路由
func (this *BaseRouter) RemoveChild(name string) (Router, bool) {
	var r, ok = this.children[name]
	if ok {
		delete(this.children, name)
		if r.Normal() {
			delete(this.normalChildren, this.unify(r.MatchString()))
		} else {
			delete(this.abnormalChildren, name)
		}
		return r, ok
	}
	return nil, false
}

// Match 匹配指定路由上下文,匹配成功则返回RouterExcutor
func (this *BaseRouter) Match(context RouterContext) (RouterExcutor, bool) {
	var segs = context.Segments()
	var data [][]string
	if !this.Normal() {
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
		var match = this.unify(segs[0])
		var r, ok = this.normalChildren[match]
		if ok {
			executor, ok = r.Match(context)
		}
		if !ok {
			for _, v := range this.abnormalChildren {
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
