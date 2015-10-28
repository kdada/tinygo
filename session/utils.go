package session

import (
	"encoding/base64"
	"math/rand"
	"net"
	"strconv"
	"strings"
	"time"
	"unsafe"
)

// guid 生成全局唯一标志
func Guid() string {
	var iface, err = net.InterfaceByIndex(0)
	var macAddr = ""
	if err == nil {
		macAddr = iface.HardwareAddr.String()
	}
	var x = new(byte)
	var seed = time.Now().UnixNano()
	var str = macAddr + strconv.Itoa(int(seed)) + strconv.Itoa(int(uintptr(unsafe.Pointer(x))))
	var data = []byte(str)
	rand.Seed(seed)
	var r = rand.Intn(len(data))
	var lendata = len(data)
	for i := 0; i < lendata; i++ {
		var temp = data[i]
		data[i] = data[r]
		data[r] = temp
		r++
		r %= lendata
	}
	var result = make([]byte, base64.StdEncoding.EncodedLen(lendata))
	base64.StdEncoding.Encode(result, data)
	return strings.ToUpper(strings.TrimRight(string(result), "="))
}
