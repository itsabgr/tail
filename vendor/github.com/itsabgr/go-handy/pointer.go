package handy

import "unsafe"

//Ref returns a uintptr points to the value
func Ref(value interface{}) uintptr {
	return uintptr(unsafe.Pointer(&value))
}

//DeRef returns the value pointer points to
func DeRef(ptr uintptr) interface{} {
	return *((*interface{})(unsafe.Pointer(ptr)))
}
