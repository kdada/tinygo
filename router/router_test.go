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

func (this *TestBaseExcutor) Execute() error {
	this.result = 124
	return nil
}

type TestFileExcutor struct {
	BaseRouterExecutor
	result int
}

func (this *TestFileExcutor) Execute() error {
	this.result = 34454
	return nil
}

type TestContext struct {
	BaseContext
	values map[string]string
}

// Value 返回路由值
func (this *TestContext) Value(name string) string {
	return this.values[name]
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
	if !ok {
		t.Fatal("基础路由查询失败")
	}
	var err = e.Execute()
	if err != nil || e.(*TestBaseExcutor).result != 124 {
		t.Fatal("基础路由执行错误")
	}
	context = new(TestContext)
	context.Segs = []string{"", "upload", "some.file"}
	e, ok = root.Match(context)
	if !ok {
		t.Fatal("无限路由查询失败")
	}
	err = e.Execute()
	if err != nil || e.(*TestFileExcutor).result != 34454 {
		t.Fatal("无限路由执行错误")
	}
}
