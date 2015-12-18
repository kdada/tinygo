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
	var hws, _ = net.Interfaces()
	var macAddr = "6a:64:77:6c:6b:6a"
	for _, h := range hws {
		if !strings.HasPrefix(h.HardwareAddr.String(), "00:00:00:00:00:00") {
			macAddr = h.HardwareAddr.String()
			break
		}
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
