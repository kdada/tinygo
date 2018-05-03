package session

import (
	"testing"

	"github.com/kdada/tinygo/util"
)

func Suite(t *testing.T, impl Session) {
	impl.SetString("sValue", "some string data.")
	var s, succ = impl.String("sValue")
	if !succ {
		t.Errorf("get string value %s failed.", "sValue")
	}
	if s != "some string data." {
		t.Errorf("get wrong string value, expected %s got %s.", "some string data.", s)
	}
	s2, succ := impl.String("sValueNotExist")
	if succ {
		t.Errorf("get string not exist: %s", s2)
	}
}

func SuiteParallel(t *testing.T, impl Session) {
	t.Run("writing1", func(t *testing.T) {
		t.Parallel()
		for i := 0; i < 10000; i++ {
			impl.SetValue("key", "value1")
		}
	})

	t.Run("reading1", func(t *testing.T) {
		t.Parallel()
		for i := 0; i < 10000; i++ {
			impl.Value("key")
		}
	})

	t.Run("writing2", func(t *testing.T) {
		t.Parallel()
		for i := 0; i < 10000; i++ {
			impl.SetValue("key", "value2")
		}
	})

	t.Run("reading2", func(t *testing.T) {
		t.Parallel()
		for i := 0; i < 10000; i++ {
			impl.Value("key")
		}
	})

	t.Run("writing3", func(t *testing.T) {
		t.Parallel()
		for i := 0; i < 10000; i++ {
			impl.SetValue("key", "value3")
		}
	})

	t.Run("reading3", func(t *testing.T) {
		t.Parallel()
		for i := 0; i < 10000; i++ {
			impl.Value("key")
		}
	})

	t.Run("writing4", func(t *testing.T) {
		t.Parallel()
		for i := 0; i < 10000; i++ {
			impl.SetValue("key", "value4")
		}
	})

	t.Run("reading4", func(t *testing.T) {
		t.Parallel()
		for i := 0; i < 10000; i++ {
			impl.Value("key")
		}
	})
}

func SuiteBenchmark(b *testing.B, impl Session) {
	for i := 0; i < b.N; i++ {
		var uuid = util.NewUUID()
		impl.SetValue(uuid.String(), uuid.Hex())
		var value, ok = impl.Value(uuid.String())
		if !ok {
			b.Error("error read key")
		}
		if value != uuid.Hex() {
			b.Errorf("return wrong value %s, expected %s", value, uuid.Hex())
		}
	}
}

func SuiteBenchmarkParallel(b *testing.B, impl Session) {
	var key = "testkey"
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			var uuid = util.NewUUID()
			impl.SetValue(key, uuid.Hex())
			var _, ok = impl.Value(key)
			if !ok {
				b.Error("error read key")
			}
		}
	})
}
