package router

import (
	"fmt"
	"testing"
)

func TestParseReg(t *testing.T) {
	var seg = "list(page)_(number=[0-9]+)"
	var segment, err = ParseReg(seg)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(segment.Exp, segment.Keys)
	var data, err2 = segment.Parse("listxxx_324234")
	if err2 != nil {
		t.Fatal(err2)
	}
	fmt.Println("data:", data)
}
