package meta

import (
	"errors"
	"fmt"
	"reflect"
	"testing"
)

type AService struct {
	kk int
}

func (this *AService) PP() string {
	return "test-test"
}

type BService struct {
	uu float32
}
type TestS struct {
	Aa string `!;clen>10`
	Bb int    `?;>100`
	Cc float32
}

type TestSS struct {
	TestS
	Value  *AService
	Common *BService
}

func Watch(t TestSS, Value *AService, param struct {
	Common *BService
}) error {
	if t.Aa != "asdsddadsad" {
		return errors.New("Aa字段值注入错误:" + t.Aa)
	}
	if t.Bb != 200 {
		return errors.New("Bb字段值注入错误:" + fmt.Sprint(t.Bb))
	}
	if Value == nil || Value.kk != 12435 || Value.PP() != "test-test" {
		return errors.New("函数参数Value注入错误:" + fmt.Sprint(Value))
	}
	if param.Common == nil || param.Common.uu-1243.522 >= 0.000001 {
		return errors.New("函数参数param.Common注入错误:" + fmt.Sprint(param.Common))
	}
	if t.Value == nil || t.Value.kk != 12435 || t.Value.PP() != "test-test" {
		return errors.New("函数参数t.Value注入错误:" + fmt.Sprint(Value))
	}
	if t.Common == nil || t.Common.uu-1243.522 >= 0.000001 {
		return errors.New("函数参数t.Common注入错误:" + fmt.Sprint(param))
	}
	return errors.New("没有错误")
}

func TestMeta(t *testing.T) {
	//注册类型和值
	GlobalValueContainer.Register(nil, func() *AService {
		return &AService{
			12435,
		}
	})
	GlobalValueContainer.Register("Common", func() *BService {
		return &BService{
			1243.522,
		}
	})
	GlobalValueContainer.Register("Aa", "asdsddadsad")
	GlobalValueContainer.Register("Bb", 200)
	// 测试方法注入
	var f = reflect.ValueOf(Watch)
	var g, e = AnalyzeMethod("", &f)
	if e != nil {
		t.Fatal(e)
	}
	var s, err = g.Generate(GlobalValueContainer)
	if err != nil {
		t.Fatal(err)
	}
	var v, ok = s.([]interface{})
	if !ok || len(v) != 1 {
		t.Fatal("注入函数返回值错误")
	}
	var r1 = v[0]
	err, ok = r1.(error)
	if !ok {
		t.Fatal("注入函数返回值类型错误")
	}
	if err.Error() != "没有错误" {
		t.Fatal(err)
	}

}
