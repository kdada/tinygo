package validator

import (
	"strconv"
)

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
