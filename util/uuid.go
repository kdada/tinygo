package util

import (
	"fmt"
	"math/rand"
	"net"
	"strings"
	"sync/atomic"
	"time"
)

// UUID 全局唯一标识符
type UUID struct {
	TimeLow                 uint32
	TimeMid                 uint16
	TimeHighAndVersion      uint16
	ClockSeqHighAndReserved uint8
	ClockSeqLow             uint8
	Node                    [6]byte
}

// 时钟序列
var clockSeq uint32 = 0

// mac地址
var macAddr []byte

// 时间偏移(100ns) 自1582-10-15
const timeOffset int64 = 122192928000000000

// UUID                   = time-low "-" time-mid "-"
//                               time-high-and-version "-"
//                               clock-seq-and-reserved
//                               clock-seq-low "-" node

// NewUUID 创建新的UUID
func NewUUID() *UUID {
	var uuid = new(UUID)
	copy(uuid.Node[:], macAddr)

	var newClockSeq = atomic.AddUint32(&clockSeq, 1)
	uuid.ClockSeqLow = byte(newClockSeq)
	uuid.ClockSeqHighAndReserved = byte(newClockSeq>>8&0x3f) | 0x80

	var nsFrom15821015 uint64 = uint64(time.Now().UnixNano()/100 + timeOffset)
	nsFrom15821015 = (nsFrom15821015 & 0x0fffffffffffffff) | 0x1000000000000000

	uuid.TimeLow = uint32(nsFrom15821015)
	uuid.TimeMid = uint16(nsFrom15821015 >> 32)
	uuid.TimeHighAndVersion = uint16(nsFrom15821015 >> 48)
	return uuid
}

// Hex 返回Hex字符串
func (this *UUID) Hex() string {
	return fmt.Sprintf("%08x%04x%04x%02x%02x%x", this.TimeLow, this.TimeMid, this.TimeHighAndVersion, this.ClockSeqHighAndReserved, this.ClockSeqLow, this.Node)
}

// String 返回UUID标准格式
func (this *UUID) String() string {
	return fmt.Sprintf("%08x-%04x-%04x-%02x%02x-%x", this.TimeLow, this.TimeMid, this.TimeHighAndVersion, this.ClockSeqHighAndReserved, this.ClockSeqLow, this.Node)
}

// init 初始化uuid变量
func init() {
	var hws, _ = net.Interfaces()
	for _, h := range hws {
		var addr = h.HardwareAddr.String()
		if addr != "" && !strings.HasPrefix(addr, "00:00:00:00:00:00") {
			macAddr = h.HardwareAddr
			break
		}
	}
	if macAddr == nil {
		panic(ErrorMacAddrNotFound)
	}
	rand.Seed(time.Now().UnixNano())
	clockSeq = uint32(rand.Int31())
}
