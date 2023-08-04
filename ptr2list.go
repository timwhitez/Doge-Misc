


func Ptr2List(ptr uintptr) []interface{} {
	return (*[0xffff]interface{})(unsafe.Pointer(ptr))[:]
}
