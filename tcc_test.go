package tcc

import (
	"strings"
	"testing"
	"unsafe"
)

func TestValue(t *testing.T) {
	s := New()
	defer s.Release()

	if err := s.Compile(`int x = 1;`); err != nil {
		t.Fatal(err)
	}
	p, err := s.Symbol("x")
	if err != nil {
		t.Fatal(err)
	}
	if expected, got := int32(1), *(*int32)(unsafe.Pointer(p)); expected != got {
		t.Errorf("expected %v, got %v", expected, got)
	}
}

func TestFunc(t *testing.T) {
	s := New()
	defer s.Release()

	if err := s.Compile(`
int uppercase(char* s) {
	int count = 0;
	while (*s) {
		if (*s >= 'a' && *s <= 'z') {
			*s -= 32;
		}
		count++;
		s++;
	}
	return count;
}`); err != nil {
		t.Fatal(err)
	}
	p, err := s.Symbol("uppercase")
	if err != nil {
		t.Fatal(err)
	}
	lower := "hello world!"
	upper := strings.ToUpper(lower)
	if s, n := callStrFunc(p, lower); s != upper || n != len(lower) {
		t.Errorf("expected: %v %v, got: %v %v", upper, len(lower), s, n)
	}
}
