package web

import (
	"reflect"
	"strings"

	"github.com/kdada/tinygo/router"
)

// 创建适用于Web App的根路由
func NewRootRouter() router.Router {
	var r, err = router.NewRouter("base", "", nil)
	if err != nil {
		panic(err)
	}
	return r
}

// 创建控制器路由
//  instance:控制器对象
//  控制器方法必须满足如下格式:
//   func (this *SomeController) Method(param *ParamStruct) web.Result
//   this:必须是控制器指针
//   param:可以没有或者有多个,如果有则类型必须为结构体指针类型
//   第一个返回结果必须能够赋值给web.Result接口
//  info:创建过程中用于输出信息的字符串
//  return:执行成功则返回router.Router
func NewControllerRouter(instance interface{}, info *string) router.Router {
	var instanceType = reflect.TypeOf(instance)
	if IsStructPtrType(instanceType) {
		panic(ErrorNotStructPtr.Format(instanceType.String()).Error())
	}
	var start *MethodMetadata = nil
	var end *MethodMetadata = nil
	var methods = make([]*MethodMetadata, 0)
	//遍历控制器方法
	ForeachMethod(instanceType, func(method reflect.Method) {
		var mMd, err = AnalyzeControllerMethod(method)
		if err != nil && info != nil {
			*info += err.Error() + "\n"
			return
		}
		// Start 和 End 方法忽略返回值检查,其他方法中符合要求的方法才能作为接口使用
		if (mMd.Name != "Start") && (mMd.Name != "End") {
			err = CheckResult(mMd)
			if err != nil && info != nil {
				*info += err.Error() + "\n"
				return
			}
		}
		switch mMd.Name {
		case "Start":
			start = mMd
		case "End":
			end = mMd
		default:
			methods = append(methods, mMd)
		}
	})
	var controllerRouter, err = router.NewRouter("base", instanceType.Name(), nil)
	if err != nil {
		panic(err)
	}
	for _, m := range methods {
		var mr, err2 = router.NewRouter("base", instanceType.Name(), nil)
		if err2 != nil {
			panic(err2)
		}
		var excutor = NewAdvancedExecutor(start, m, end)
		mr.AddChildren(HttpResultRouter(m.Return[0].Name(), func() router.RouterExcutor {
			return excutor
		}))
		controllerRouter.AddChild(mr)
	}
	return controllerRouter
}

// 创建函数路由
//  name:路由名称
//  function:函数
//  函数必须满足如下格式:
//   func Method(param *ParamStruct) web.Result
//   param:可以没有或者有多个,如果有则类型必须为结构体指针类型
//   第一个返回结果必须能够赋值给web.Result接口
//  return:执行成功则返回router.Router
func NewFuncRouter(name string, function interface{}) router.Router {
	var mMd, err = AnalyzeMethod(name, reflect.ValueOf(function))
	if err != nil {
		panic(err)
	}
	var mr, err2 = router.NewRouter("base", name, name)
	if err2 != nil {
		panic(err2)
	}
	var excutor = NewAdvancedExecutor(nil, mMd, nil)
	var mName = ""
	if CheckResult(mMd) == nil {
		mName = mMd.Return[0].String()
	}
	mr.AddChildren(HttpResultRouter(mName, func() router.RouterExcutor {
		return excutor
	}))
	return mr
}

// 创建函数路由,可匹配无限层级
//  name:路由名称
//  function:函数
//  函数必须满足如下格式:
//   func Method(param *ParamStruct) web.Result
//   param:可以没有或者有多个,如果有则类型必须为结构体指针类型
//   第一个返回结果必须能够赋值给web.Result接口
//  return:执行成功则返回router.Router
func NewMutableFuncRouter(name string, function interface{}) router.Router {
	var mMd, err = AnalyzeMethod(name, reflect.ValueOf(function))
	if err != nil {
		panic(err)
	}
	var mr, err2 = router.NewRouter("unlimited", name, name)
	if err2 != nil {
		panic(err2)
	}
	var excutor = NewAdvancedExecutor(nil, mMd, nil)
	mr.SetRouterExcutorGenerator(func() router.RouterExcutor {
		return excutor
	})
	return mr
}

// 创建空间路由
//  name:路由名称
//  return:执行成功则返回router.Router
func NewSpaceRouter(name string) router.Router {
	var r, err = router.NewRouter("base", name, name)
	if err != nil {
		panic(err)
	}
	return r
}

// 检查元数据的第一个返回值是否符合web.Result接口
func CheckResult(m *MethodMetadata) error {
	if len(m.Return) <= 0 {
		return ErrorNoReturn.Format(m.Name).Error()
	}
	var resultType = reflect.TypeOf((*Result)(nil)).Elem()
	if !m.Return[0].AssignableTo(resultType) {
		return ErrorFirstReturnMustBeResult.Format(m.Return[0]).Error()
	}
	return nil
}

// 提取name中包含的Http方法名,如果不包含任何方法名,则返回Post
func HttpMethod(name string) []string {
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

// 生成Http方法名路由
func HttpResultRouter(resultName string, gen router.RouterExcutorGenerator) []router.Router {
	var rs = make([]router.Router, 0, 1)
	for _, httpMethod := range HttpMethod(resultName) {
		var hr, err3 = router.NewRouter("base", httpMethod, nil)
		if err3 != nil {
			panic(err3)
		}
		hr.SetRouterExcutorGenerator(gen)
		rs = append(rs, hr)
	}
	return rs
}
