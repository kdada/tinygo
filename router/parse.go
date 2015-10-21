package router

import (
	"regexp"
	"strings"
)

//正则段
type RegSegment struct {
	Exp    string         //正则表达式
	Regexp *regexp.Regexp //编译后的正则表达式
	Keys   []string       //可提取keys
}

// Parse 解析字符串并生成key-value形式的值
func (this *RegSegment) Parse(str string) (map[string]string, error) {
	var data = this.Regexp.FindAllStringSubmatch(str, 1)
	if data != nil && len(data) >= 1 {
		var result = data[0][1:]
		if len(result) == len(this.Keys) {
			var kvs = make(map[string]string, len(this.Keys))
			for i, v := range result {
				kvs[this.Keys[i]] = v
			}
			return kvs, nil
		}
	}
	return nil, RouterErrorRegexpNotMatchError.Error()
}

// ParseReg 解析正则路由段字符串
// return:如果可以解析出正则内容,则返回RegSegment,否则返回nil
func ParseReg(exp string) (*RegSegment, error) {
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
					key = reg[:pos]
					value = reg[pos+1:]
				} else if pos == -1 {
					key = reg
					value = ".*"
				} else {
					//错误
					return nil, RouterErrorRegexpFormatError.Format(reg).Error()
				}
				if key != "" {
					rs.Keys = append(rs.Keys, key)
					rs.Exp += "(" + value + ")"
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
		var err error
		rs.Regexp, err = regexp.Compile(rs.Exp)
		if err != nil {
			return nil, RouterErrorRegexpParseError.Format(rs.Exp).Error()
		}
		return rs, nil
	}
	return nil, RouterErrorRegexpNoneError.Format(exp).Error()
}
