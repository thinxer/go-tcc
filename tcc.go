package tcc

/*
#cgo LDFLAGS: -ltcc
#include <libtcc.h>
#include <stdlib.h>
*/
import "C"
import (
	"errors"
	"unsafe"
)

type State struct {
	h *C.TCCState

	relocated bool
}

func New() *State {
	h := C.tcc_new()
	return &State{h: h}
}

func (s *State) Release() {
	C.tcc_delete(s.h)
}

func (s *State) Compile(prog string) error {
	cprog := C.CString(prog)
	defer C.free(unsafe.Pointer(cprog))
	if ret := C.tcc_compile_string(s.h, cprog); ret < 0 {
		return errors.New("compilation failed")
	}
	return nil
}

func (s *State) Symbol(name string) (uintptr, error) {
	if !s.relocated {
		// TODO: manually manage memory
		if ret := C.tcc_relocate(s.h, unsafe.Pointer(uintptr(1))); ret < 0 {
			return 0, errors.New("relocate error")
		}
		s.relocated = true
	}

	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))

	if h := uintptr(unsafe.Pointer(C.tcc_get_symbol(s.h, cname))); h == 0 {
		return 0, errors.New("not found")
	} else {
		return h, nil
	}
}
