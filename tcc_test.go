package tcc

import (
	"testing"
	"unsafe"
)

func TestValue(t *testing.T) {
	s := New()
	s.Compile(`
        int x = 1;
    `)
	p, err := s.Symbol("x")
	if err != nil {
		t.Fatal(err)
	}
	if expected, got := int32(1), *(*int32)(unsafe.Pointer(p)); expected != got {
		t.Errorf("expected %v, got %v", expected, got)
	}
}
