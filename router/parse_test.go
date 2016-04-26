package router

import "testing"

func TestParseReg(t *testing.T) {
	var seg = "list{page}_{number=[0-9]+}"
	var segment, err = ParseReg(seg)
	if err != nil || segment.Exp != "^list(.*)_([0-9]+)$" || segment.Keys[0] != "page" || segment.Keys[1] != "number" {
		t.Fatal(err, "解析错误")
	}
	var data, err2 = segment.Parse("listxxx_324234")
	if err2 != nil || data["page"] != "xxx" || data["number"] != "324234" {
		t.Fatal(err2, "匹配数据错误")
	}
}
