package validator

import (
	"testing"
)

func TestStringValidator(t *testing.T) {
	var count = 0
	var src = ` (<= 3000 && >-213) || /[1]/ && aboveNow('2012-12-12')`
	RegisterFunc("aboveNow", func(str string, sss string) bool {
		if sss != "2012-12-12" {
			t.Fatal("参数错误")
		}
		count++
		return true
	})
	var v, err = NewValidator("string", src)
	if err != nil {
		t.Fatal(err)
	}
	if v.Validate("-213") != true {
		t.Fatal("校验失败 -213")
	}
	if v.Validate("3") != true {
		t.Fatal("校验失败 3")
	}
	if v.Validate("324") != true {
		t.Fatal("校验失败 324")
	}
	if v.Validate("3245") != false {
		t.Fatal("校验失败 3245")
	}
	if v.Validate("54545") != false {
		t.Fatal("校验失败 54545")
	}
	if count != 3 {
		t.Fatal("校验失败")
	}
}
