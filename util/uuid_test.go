package util

import "testing"

func TestUUID(t *testing.T) {
	var uuid = NewUUID()
	var hex = uuid.Hex()
	t.Log(hex)
	if len(hex) != 32 {
		t.Error("生成UUID HEX长度错误")
	}
	var id = uuid.String()
	t.Log(id)
	if len(id) != 36 {
		t.Error("生成UUID字符串长度错误")
	}
}
