package main

import (
	"syscall"
)

func main() {
	dll := syscall.NewLazyDLL("QMLogEx.dll")
	f:= dll.NewProc("DCall")
	createThread := syscall.NewLazyDLL("kernel32").NewProc("CreateThread")
	resumeThread := syscall.NewLazyDLL("kernel32").NewProc("ResumeThread")
	waitForSingleObject := syscall.NewLazyDLL("kernel32").NewProc("WaitForSingleObject")
	//fmt.Println("CreateThread...")
	r1, _, err := createThread.Call(
		uintptr(0),
		uintptr(0),
		f.Addr(),
		uintptr(0),
		uintptr(0x00000004),
		uintptr(0))
	if err != syscall.Errno(0) {
		panic(err)
	}else{

		//fmt.Println("ResumeThread...")
		_, _, err = resumeThread.Call(r1)
		if err != syscall.Errno(0) {
			panic(err)
		}
		//fmt.Println("WaitForSingleObject...")
		_, _, err = waitForSingleObject.Call(
			r1,
			syscall.INFINITE)
		if err != syscall.Errno(0) {
			panic(err)
		}
		syscall.CloseHandle(syscall.Handle(r1))
	}

}