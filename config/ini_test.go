package config

import "testing"

func TestIniConfig(t *testing.T) {
	var config, err = NewConfig(ConfigTypeIni, "test_bom.cfg")
	if err != nil {
		t.Error(err)
		return
	}
	var global = config.GlobalSection()
	var configk, _ = global.String("config")
	if configk != "test" {
		t.Error("config(键)错误")
	}
	var ckey, _ = global.String("中文Key")
	if ckey != "中文测试" {
		t.Error("中文Key(键)错误")
	}
	var debug, _ = config.Section("debug")
	var dt, _ = debug.Bool("Debug-Test")
	if !dt {
		t.Error("Debug-Test(键)错误")
	}
	var cmd, _ = debug.Int("Command")
	if cmd != 123465 {
		t.Error("Command(键)错误")
	}
	var fl, _ = debug.Float("DebugFloat")
	if fl != 12.11121 {
		t.Error("DebugFloat(键)错误")
	}
	var st, _ = debug.String("TestString")
	if st != "asdasd哈撒地方阿斯蒂芬asdas'\\[]..///.'']]ppp'" {
		t.Error("TestString(键)错误")
	}
	var release, _ = config.Section("release")
	t.Log(release)
}
