package session

import (
	"testing"

	"github.com/kdada/tinygo/util"
)

func TestMemSession(t *testing.T) {
	var ssid = util.NewUUID().Hex()
	var ss = newMemSession(ssid)
	var ssid2 = ss.SessionId()
	if ssid != ssid2 {
		t.Errorf("ssid error, expected %s, got %s", ssid, ssid2)
	}
	Suite(t, ss)
}

func TestMemSessionParallel(t *testing.T) {
	var ssid = util.NewUUID().Hex()
	var ss = newMemSession(ssid)
	SuiteParallel(t, ss)
}

func BenchmarkMemSession(b *testing.B) {
	var ssid = util.NewUUID().Hex()
	var ss = newMemSession(ssid)
	SuiteBenchmark(b, ss)
}

func BenchmarkMemSessionParallel(b *testing.B) {
	var ssid = util.NewUUID().Hex()
	var ss = newMemSession(ssid)
	SuiteBenchmarkParallel(b, ss)
}
