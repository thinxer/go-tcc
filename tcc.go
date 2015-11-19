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

	"github.com/edsrzf/mmap-go"
)

var ErrCompile = errors.New("compilation error")
var ErrRelocate = errors.New("relocate error")
var ErrSymbolNotFound = errors.New("symbol not found")

type State struct {
	h *C.TCCState

	relocated bool
	code      mmap.MMap
}

func New() *State {
	h := C.tcc_new()
	return &State{h: h}
}

func (s *State) Release() {
	C.tcc_delete(s.h)
	if s.relocated {
		s.code.Unmap()
	}
}

func (s *State) Compile(prog string) error {
	cprog := C.CString(prog)
	defer C.free(unsafe.Pointer(cprog))
	if ret := C.tcc_compile_string(s.h, cprog); ret < 0 {
		return ErrCompile
	}
	return nil
}

func (s *State) relocate() (err error) {
	if s.relocated {
		return nil
	}

	size := C.tcc_relocate(s.h, unsafe.Pointer(uintptr(0)))
	if size < 0 {
		return ErrRelocate
	}
	s.code, err = mmap.MapRegion(nil, int(size), mmap.EXEC|mmap.RDWR, mmap.ANON, int64(0))
	if err != nil {
		return
	}
	if ret := C.tcc_relocate(s.h, unsafe.Pointer(&s.code[0])); ret < 0 {
		s.code.Unmap()
		return ErrRelocate
	}
	s.relocated = true

	return nil
}

func (s *State) Symbol(name string) (uintptr, error) {
	if err := s.relocate(); err != nil {
		return 0, err
	}

	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))

	if h := uintptr(unsafe.Pointer(C.tcc_get_symbol(s.h, cname))); h == 0 {
		return 0, ErrSymbolNotFound
	} else {
		return h, nil
	}
}
