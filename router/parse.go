package router

import (
	"strings"
)

//正则段
type RegSegment struct {
	Exp  string   //正则表达式
	Keys []string //可提取keys
}

// ParseReg 解析正则路由段字符串
func ParseReg(exp string) *RegSegment {
	var rs = new(RegSegment)
	rs.Keys = make([]string, 0)
	var bytes = []byte(exp)

	var lastSegStart = 0
	for i, b := range bytes {
		switch b {
		case 123: //"{"
			{
				//截取普通字符串
				rs.Exp += string(bytes[lastSegStart:i])
				lastSegStart = i + 1
			}
		case 125: //"}"
			{
				//截取正则字符串
				var reg = string(bytes[lastSegStart:i])
				var pos = strings.Index(reg, "=")
				var key, value = "", ""
				if pos > 0 {
					pos += lastSegStart
					key = string(bytes[lastSegStart:pos])
					value = string(bytes[pos+1 : i])
				}
				if key != "" {
					rs.Keys = append(rs.Keys, key)
					if value == "" {
						value = ".*"
					}
					rs.Exp += value
				}
				lastSegStart = i + 1
			}
		}
	}
	if len(rs.Keys) > 0 {
		if lastSegStart < len(bytes)-1 {
			rs.Exp += string(bytes[lastSegStart:len(bytes)])
		}
		rs.Exp = "^" + rs.Exp + "$"
		return rs
	}
	return nil
}

// RouterSegment 将整个路由字符串分段
func RouterSegment(router string) []string {
	return nil
}
