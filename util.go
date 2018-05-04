package onering

import "unsafe"

type iface struct {
	t, d unsafe.Pointer
}

func extractptr(i interface{}) unsafe.Pointer {
	if DEBUG {
		checkPutType(i, getCallerPath())
	}
	return (*iface)(unsafe.Pointer(&i)).d
}

func extractfn(i interface{}) func(Iter, unsafe.Pointer) {
	if DEBUG {
		checkFuncType(i, getCallerPath())
	}
	var ptr = (*iface)(unsafe.Pointer(&i)).d
	return *(*func(Iter, unsafe.Pointer))(unsafe.Pointer(&ptr))
}

func inject(i interface{}, ptr unsafe.Pointer) {
	if DEBUG {
		checkGetType(i, getCallerPath())
	}
	var v = (*unsafe.Pointer)((*iface)(unsafe.Pointer(&i)).d)
	*v = ptr
}
