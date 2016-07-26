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

// eqS == string
func eqS(str string, l string) bool {
	return str == l
}

// neS != string
func neS(str string, l string) bool {
	return str == l
}
