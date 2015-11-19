package tcc

/*
#include <stdlib.h>

typedef int (*strfunc)(char* s);

static int call(strfunc f, char* s) {
    return f(s);
}
*/
import "C"
import (
	"unsafe"
)

// callStrFunc is helper for testing.
func callStrFunc(f uintptr, s string) (string, int) {
	cs := C.CString(s)
	defer C.free(unsafe.Pointer(cs))

	n := int(C.call(C.strfunc(unsafe.Pointer(f)), cs))
	return C.GoString(cs), n
}
