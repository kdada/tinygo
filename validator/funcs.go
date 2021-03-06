package validator

import (
	"strconv"
	"unicode/utf8"
)

// 整数比较方法
// ltI < int64
func ltI(str string, l int64) bool {
	var i, err = strconv.ParseInt(str, 10, 64)
	return err == nil && i < l
}

// leI <= int64
func leI(str string, l int64) bool {
	var i, err = strconv.ParseInt(str, 10, 64)
	return err == nil && i <= l
}

// gtI > int64
func gtI(str string, l int64) bool {
	var i, err = strconv.ParseInt(str, 10, 64)
	return err == nil && i > l
}

// geI >= int64
func geI(str string, l int64) bool {
	var i, err = strconv.ParseInt(str, 10, 64)
	return err == nil && i >= l
}

// eqI == int64
func eqI(str string, l int64) bool {
	var i, err = strconv.ParseInt(str, 10, 64)
	return err == nil && i == l
}

// neI != int64
func neI(str string, l int64) bool {
	var i, err = strconv.ParseInt(str, 10, 64)
	return err == nil && i != l
}

// 浮点数比较方法
// ltF < float64
func ltF(str string, l float64) bool {
	var i, err = strconv.ParseFloat(str, 64)
	return err == nil && i < l
}

// leF <= float64
func leF(str string, l float64) bool {
	var i, err = strconv.ParseFloat(str, 64)
	return err == nil && i <= l
}

// gtF > float64
func gtF(str string, l float64) bool {
	var i, err = strconv.ParseFloat(str, 64)
	return err == nil && i > l
}

// geF >= float64
func geF(str string, l float64) bool {
	var i, err = strconv.ParseFloat(str, 64)
	return err == nil && i >= l
}

// eqF == float64
func eqF(str string, l float64) bool {
	var i, err = strconv.ParseFloat(str, 64)
	return err == nil && i == l
}

// neF != float64
func neF(str string, l float64) bool {
	var i, err = strconv.ParseFloat(str, 64)
	return err == nil && i != l
}

//字符串比较方法
// eqS == string
func eqS(str string, l string) bool {
	return str == l
}

// neS != string
func neS(str string, l string) bool {
	return str == l
}

// 字节数比较方法
// lenLtI < int64
func lenLtI(str string, l int64) bool {
	return len(str) < int(l)
}

// lenLeI <= int64
func lenLeI(str string, l int64) bool {
	return len(str) <= int(l)
}

// lenGtI > int64
func lenGtI(str string, l int64) bool {
	return len(str) > int(l)
}

// lenGeI >= int64
func lenGeI(str string, l int64) bool {
	return len(str) >= int(l)
}

// lenEqI == int64
func lenEqI(str string, l int64) bool {
	return len(str) == int(l)
}

// lenNeI != int64
func lenNeI(str string, l int64) bool {
	return len(str) != int(l)
}

// 字符数比较方法
// clenLtI < int64
func clenLtI(str string, l int64) bool {
	return utf8.RuneCountInString(str) < int(l)
}

// clenLeI <= int64
func clenLeI(str string, l int64) bool {
	return utf8.RuneCountInString(str) <= int(l)
}

// clenGtI > int64
func clenGtI(str string, l int64) bool {
	return utf8.RuneCountInString(str) > int(l)
}

// clenGeI >= int64
func clenGeI(str string, l int64) bool {
	return utf8.RuneCountInString(str) >= int(l)
}

// clenEqI == int64
func clenEqI(str string, l int64) bool {
	return utf8.RuneCountInString(str) == int(l)
}

// clenNeI != int64
func clenNeI(str string, l int64) bool {
	return utf8.RuneCountInString(str) != int(l)
}
