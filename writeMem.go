package main

import "unsafe"

func main(){
	var shellcode = []byte{0x00}
	var baseA uintptr


	//write memory
	WriteMemory(shellcode, baseA)
	memcpy(baseA, shellcode)
	//	*(**uintptr)(unsafe.Pointer(&targetPtr)) = (*uintptr)(unsafe.Pointer(&rawShellcode))
	//syscall.Syscall(uintptr(unsafe.Pointer(&rawShellcode[0])),0, 0, 0, 0)
	//targetPtr()
}


//WriteMemory writes the provided memory to the specified memory address. Does **not** check permissions, may cause panic if memory is not writable etc.
func WriteMemory(inbuf []byte, destination uintptr) {
	for index := uint32(0); index < uint32(len(inbuf)); index++ {
		writePtr := unsafe.Pointer(destination + uintptr(index))
		v := (*byte)(writePtr)
		*v = inbuf[index]
	}
}

func memcpy(base uintptr, buf []byte) {
	for i := 0; i < len(buf); i++ {
		*(*byte)(unsafe.Pointer(base + uintptr(i))) = buf[i]
	}
}

func memcpy(src, dst, size uintptr) {
	for i := uintptr(0); i < size; i++ {
		*(*uint8)(unsafe.Pointer(dst + i)) = *(*uint8)(unsafe.Pointer(src + i))
	}
}

func Memset(ptr uintptr, c byte, n uintptr){
	var i uintptr
	for i = 0;i<n;i++{
		pByte:=(*byte)(unsafe.Pointer(ptr+i))
		*pByte = c
	}
}
