package web

import (
	"reflect"
	"strings"

	"github.com/kdada/tinygo/meta"
	"github.com/kdada/tinygo/router"
)

// Http方法
type HttpMethod string

const (
	HttpMethodGet     HttpMethod = "Get"     // Get方法
	HttpMethodPost    HttpMethod = "Post"    // Post方法
	HttpMethodPut     HttpMethod = "Put"     // Put方法
	HttpMethodDelete  HttpMethod = "Delete"  // Delete方法
	HttpMethodOptions HttpMethod = "Options" // Options方法
	HttpMethodHead    HttpMethod = "Head"    // Head方法
	HttpMethodTrace   HttpMethod = "Trace"   // Trace方法
	HttpMethodConnect HttpMethod = "Connect" // Connect方法
)

// 路由方法信息
type RouterMethod struct {
	MethodName string     //控制器方法名称
	RouterName string     //路由名称
	HttpMethod HttpMethod //Http方法名称
}

// NewSpaceRouter 创建空间路由
//  name:路由名称
//  return:执行成功则返回router.Router
func NewSpaceRouter(name string) router.Router {
	var r, err = router.NewRouter("base", name, name)
	if err != nil {
		panic(err)
	}
	return r
}

// NewSpaceRouters 创建多级空间路由
//  url:路由,使用/分割
//  return:执行成功则返回(根路由,叶路由)
func NewSpaceRouters(url string) (router.Router, router.Router) {
	var names = strings.Split(url, "/")
	var root router.Router
	var leaf router.Router
	for _, v := range names {
		v = strings.TrimSpace(v)
		if v == "" {
			continue
		}
		var r = NewSpaceRouter(v)
		if root == nil {
			root = r
			leaf = r
			continue
		}
		leaf.AddChild(r)
		leaf = r
	}
	return root, leaf
}

// NewRootRouter 创建适用于Web App的根路由
func NewRootRouter() router.Router {
	return NewSpaceRouter("")
}

// NewControllerRouter 创建控制器路由,根据方法返回值确定该方法处理哪种形式的http请求
//  instance:控制器对象
//  控制器方法必须满足如下格式:
//   func (this *SomeController) Method(param *ParamStruct) web.Result
//   this:必须是控制器指针
//   param:可以没有或者有多个,如果有则类型必须为结构体指针类型
//   第一个返回结果最好是能够赋值给web.Result接口,也可以是其他类型
//  return:执行成功则返回控制器的router.Router
func NewControllerRouter(instance interface{}) router.Router {
	var instanceType = reflect.TypeOf(instance)
	if !meta.IsStructPtrType(instanceType) {
		panic(ErrorNotStructPtr.Format(instanceType.String()).Error())
	}
	var methods = make([]*meta.MethodMetadata, 0)
	//遍历控制器方法
	var err = meta.ForeachMethod(instanceType, func(method reflect.Method) error {
		var mMd, err = meta.AnalyzeStructMethod(&method)
		if err != nil {
			return err
		}
		// 对返回值进行检查,符合要求的方法才能作为接口使用
		err = CheckResult(mMd)
		if err != nil {
			return err
		}
		methods = append(methods, mMd)
		return nil
	})
	if err != nil {
		panic(err)
	}
	var controllerName = instanceType.Elem().Name()
	controllerName = strings.TrimSuffix(controllerName, "Controller")
	var controllerRouter = NewSpaceRouter(controllerName)
	for _, m := range methods {
		var mr = NewSpaceRouter(m.Name)
		var excutor = NewAdvancedExecutor(m)
		mr.AddChildren(HttpResultRouter(m.Return[0].Name(), func() router.RouterExcutor {
			return excutor
		}))
		controllerRouter.AddChild(mr)
	}
	return controllerRouter
}

// NewCustomControllerRouter 创建定制化控制器路由
//  instance:控制器对象
//  name:控制器路由名称,为空时使用instance类名(不含Controller)
//  methodsInfo:路由方法信息数组,RouterName为空时使用方法名
//  return:执行成功则返回控制器的router.Router
func NewCustomControllerRouter(instance interface{}, name string, methodsInfo []RouterMethod) router.Router {
	var instanceType = reflect.TypeOf(instance)
	if !meta.IsStructPtrType(instanceType) {
		panic(ErrorNotStructPtr.Format(instanceType.String()).Error())
	}
	var methods = make([]*meta.MethodMetadata, len(methodsInfo))
	//遍历方法
	for k, v := range methodsInfo {
		var method, ok = instanceType.MethodByName(v.MethodName)
		if !ok {
			panic(ErrorNoSpecificMethod.Format(instanceType.String(), v.MethodName).Error())
		}
		var mMd, err = meta.AnalyzeStructMethod(&method)
		if err != nil {
			panic(err)
		}
		methods[k] = mMd
	}
	// 生成路由
	if name == "" {
		name = strings.TrimSuffix(instanceType.Elem().Name(), "Controller")
	}
	var controllerRouter = NewSpaceRouter(name)
	for i, m := range methods {
		var info = methodsInfo[i]
		var rname = info.RouterName
		if rname == "" {
			rname = m.Name
		}
		var mr = NewSpaceRouter(rname)
		var excutor = NewAdvancedExecutor(m)
		mr.AddChildren(HttpResultRouter(string(info.HttpMethod), func() router.RouterExcutor {
			return excutor
		}))
		controllerRouter.AddChild(mr)
	}
	return controllerRouter
}

// NewFuncRouter 创建函数路由,根据方法返回值确定该方法处理哪种形式的http请求
//  name:路由名称
//  function:函数
//  函数必须满足如下格式:
//   func Method(param *ParamStruct) web.Result
//   param:可以没有或者有多个,如果有则类型必须为结构体指针类型
//   第一个返回结果最好是能够赋值给web.Result接口,也可以是其他类型
//  return:执行成功则返回router.Router
func NewFuncRouter(name string, function interface{}) router.Router {
	var v = reflect.ValueOf(function)
	var mMd, err = meta.AnalyzeMethod(name, &v)
	if err != nil {
		panic(err)
	}
	var mr = NewSpaceRouter(name)
	var excutor = NewAdvancedExecutor(mMd)
	var mName = ""
	if CheckResult(mMd) == nil {
		mName = mMd.Return[0].String()
	}
	mr.AddChildren(HttpResultRouter(mName, func() router.RouterExcutor {
		return excutor
	}))
	return mr
}

// NewMutableFuncRouter 创建函数路由,可匹配无限层级和任意http方法的请求
//  name:路由名称
//  function:函数
//  函数必须满足如下格式:
//   func Method(param *ParamStruct) web.Result
//   param:可以没有或者有多个,如果有则类型必须为结构体指针类型
//   第一个返回结果最好是能够赋值给web.Result接口,也可以是其他类型
//  return:执行成功则返回router.Router
func NewMutableFuncRouter(name string, function interface{}) router.Router {
	var v = reflect.ValueOf(function)
	var mMd, err = meta.AnalyzeMethod(name, &v)
	if err != nil {
		panic(err)
	}
	var mr, err2 = router.NewRouter("unlimited", name, name)
	if err2 != nil {
		panic(err2)
	}
	var excutor = NewAdvancedExecutor(mMd)
	mr.SetRouterExcutorGenerator(func() router.RouterExcutor {
		return excutor
	})
	return mr
}

// NewFileRouter 创建文件路由,只能匹配Get类型的文件请求,返回指定的文件
//  name:路由名称
//  path:文件路径
//  return:执行成功则返回router.Router
func NewFileRouter(name string, path string) router.Router {
	var mr = NewSpaceRouter(name)
	var excutor = NewFileExecutor(path)
	mr.AddChildren(HttpResultRouter("Get", func() router.RouterExcutor {
		return excutor
	}))
	return mr
}

// NewStaticRouter 创建静态文件路由,只能匹配Get类型的文件请求,返回指定的文件
//  name:路由名称
//  path:文件目录路径
//  return:执行成功则返回router.Router
func NewStaticRouter(name string, path string) router.Router {
	var mr, err = router.NewRouter("unlimited", name, name)
	if err != nil {
		panic(err)
	}
	var excutor = NewStaticExecutor(path)
	mr.SetRouterExcutorGenerator(func() router.RouterExcutor {
		return excutor
	})
	var pathRouter = NewSpaceRouter(name)
	pathRouter.AddChild(mr)
	return pathRouter
}

// CheckResult 检查元数据的第一个返回值是否符合web.Result接口
func CheckResult(m *meta.MethodMetadata) error {
	if len(m.Return) <= 0 {
		return ErrorNoReturn.Format(m.Name).Error()
	}
	var resultType = reflect.TypeOf((*Result)(nil)).Elem()
	if !m.Return[0].AssignableTo(resultType) {
		return ErrorFirstReturnMustBeResult.Format(m.Return[0]).Error()
	}
	return nil
}

// HttpMethod 提取name中包含的Http方法名,如果不包含任何方法名,则返回Post
func HttpMethodName(name string) []string {
	var dotPos = strings.LastIndex(name, ".")
	if dotPos > 0 {
		name = name[dotPos+1:]
	}
	var result = make([]string, 0)
	if strings.Contains(name, "Get") {
		result = append(result, "Get")
	}
	if strings.Contains(name, "Put") {
		result = append(result, "Put")
	}
	if strings.Contains(name, "Delete") {
		result = append(result, "Delete")
	}
	if strings.Contains(name, "Options") {
		result = append(result, "Options")
	}
	if strings.Contains(name, "Head") {
		result = append(result, "Head")
	}
	if strings.Contains(name, "Trace") {
		result = append(result, "Trace")
	}
	if strings.Contains(name, "Connect") {
		result = append(result, "Connect")
	}
	// 如果不指定方法则默认使用Post
	if len(result) <= 0 {
		result = append(result, "Post")
	}
	return result
}

// HttpResultRouter 生成Http方法名路由
func HttpResultRouter(resultName string, gen router.RouterExcutorGenerator) []router.Router {
	var rs = make([]router.Router, 0, 1)
	for _, httpMethod := range HttpMethodName(resultName) {
		var hr, err3 = router.NewRouter("base", httpMethod, nil)
		if err3 != nil {
			panic(err3)
		}
		hr.SetRouterExcutorGenerator(gen)
		rs = append(rs, hr)
	}
	return rs
}
