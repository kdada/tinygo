package router

import (
	"fmt"
	"os"
	"testing"
)

var root Router

type TestBaseExcutor struct {
	BaseRouterExecutor
	result int
}

func (this *TestBaseExcutor) Execute() (interface{}, error) {
	this.result = 124
	return this.result, nil
}

type TestFileExcutor struct {
	BaseRouterExecutor
	result int
}

func (this *TestFileExcutor) Execute() (interface{}, error) {
	this.result = 34454
	return this.result, nil
}

type TestContext struct {
	BaseContext
	values map[string]string
}

// Value 返回路由值
func (this *TestContext) Value(name string) (string, bool) {
	var v, ok = this.values[name]
	return v, ok
}

// SetValue 设置路由值
func (this *TestContext) SetValue(name string, value string) {
	this.values[name] = value
}

// Data 返回路由上下文携带的信息
func (this *TestContext) Data() interface{} {
	return this.values
}

func TestMain(m *testing.M) {
	var err error
	root, err = NewRouter("base", "", nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	var r, _ = NewRouter("base", "Home", nil)
	root.AddChild(r)
	var r2, _ = NewRouter("base", "Index.html", nil)
	r2.SetRouterExcutorGenerator(func() RouterExcutor {
		return new(TestBaseExcutor)
	})
	r.AddChild(r2)
	r, _ = NewRouter("unlimited", "file", nil)
	r.SetRouterExcutorGenerator(func() RouterExcutor {
		return new(TestFileExcutor)
	})
	root.AddChild(r)
	os.Exit(m.Run())
}

func TestRouter(t *testing.T) {
	var context = new(TestContext)
	context.Segs = []string{"", "hOme", "indEx.html"}
	var e, ok = root.Match(context)
	if context.Level != 3 {
		t.Fatal("匹配的路由数量错误")
	}
	if !ok {
		t.Fatal("基础路由查询失败")
	}
	var r, err = e.Execute()
	if err != nil || e.(*TestBaseExcutor).result != 124 || r.(int) != 124 {
		t.Fatal("基础路由执行错误")
	}
	context = new(TestContext)
	context.Segs = []string{"", "upload", "some.file"}
	e, ok = root.Match(context)
	if context.Level != 2 {
		t.Fatal("匹配的路由数量错误")
	}
	if !ok {
		t.Fatal("无限路由查询失败")
	}
	r, err = e.Execute()
	if err != nil || e.(*TestFileExcutor).result != 34454 || r.(int) != 34454 {
		t.Fatal("无限路由执行错误")
	}
}

func TestFind(t *testing.T) {
	// 测试查找方法
	var c = NewBaseContext("/home/")
	var r, ok = root.Find(c)
	if c.Level != 2 {
		t.Fatal("匹配的路由数量错误")
	}
	if !ok || r.Name() != "Home" {
		t.Fatal("路由查找错误")
	}
	c = NewBaseContext("/ss/dddd")
	r, ok = root.Find(c)
	if c.Level != 2 {
		t.Fatal("匹配的路由数量错误")
	}
	if !ok || r.Name() != "file" {
		t.Fatal("路由查找错误")
	}
}
