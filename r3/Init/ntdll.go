package Init

import (
	"unsafe"

	"github.com/pkg/errors"
	"golang.org/x/sys/windows"
)

// reference:
// https://github.com/gentilkiwi/mimikatz/blob/master/mimikatz/modules/kuhl_m_privilege.h
// https://github.com/gentilkiwi/mimikatz/blob/master/mimikatz/modules/kuhl_m_privilege.c

var (
	modNTDLL = windows.NewLazySystemDLL("ntdll.dll")

	procRtlAdjustPrivilege = modNTDLL.NewProc("RtlAdjustPrivilege")
)

// about privilege id
const (
	SESecurity       uint32 = 8
	SELoadDriver     uint32 = 10
	SESystemTime     uint32 = 12
	SESystemProf     uint32 = 13
	SEBackup         uint32 = 17
	SEShutdown       uint32 = 19
	SEDebug          uint32 = 20
	SESystemEnv      uint32 = 22
	SERemoteShutdown uint32 = 24
)

// RtlAdjustPrivilege is used to adjust privilege with procRtlAdjustPrivilege.
func RtlAdjustPrivilege(id uint32, enable, currentThread bool) (bool, error) {
	op := "disable"
	var p0 uint32
	if enable {
		p0 = 1
		op = "enable"
	}
	var p1 uint32
	if currentThread {
		p1 = 1
	}
	var previous bool
	ret, _, _ := procRtlAdjustPrivilege.Call(
		uintptr(id), uintptr(p0), uintptr(p1), uintptr(unsafe.Pointer(&previous)),
	) // #nosec
	if ret != 0 {
		return false, errors.Errorf("failed to %s privilege: %d, error: 0x%08X", op, id, ret)
	}
	return previous, nil
}

// RtlEnableSecurity is used to enable security privilege that call RtlAdjustPrivilege.
func RtlEnableSecurity() (bool, error) {
	return RtlAdjustPrivilege(SESecurity, true, false)
}

// RtlDisableSecurity is used to disable security privilege that call RtlAdjustPrivilege.
func RtlDisableSecurity() (bool, error) {
	return RtlAdjustPrivilege(SESecurity, false, false)
}

// RtlEnableLoadDriver is used to enable load driver privilege that call RtlAdjustPrivilege.
func RtlEnableLoadDriver() (bool, error) {
	return RtlAdjustPrivilege(SELoadDriver, true, false)
}

// RtlDisableLoadDriver is used to disable load driver privilege that call RtlAdjustPrivilege.
func RtlDisableLoadDriver() (bool, error) {
	return RtlAdjustPrivilege(SELoadDriver, false, false)
}

// RtlEnableSystemTime is used to enable system time privilege that call RtlAdjustPrivilege.
func RtlEnableSystemTime() (bool, error) {
	return RtlAdjustPrivilege(SESystemTime, true, false)
}

// RtlDisableSystemTime is used to disable system time privilege that call RtlAdjustPrivilege.
func RtlDisableSystemTime() (bool, error) {
	return RtlAdjustPrivilege(SESystemTime, false, false)
}

// RtlEnableSystemProf is used to enable system profile privilege that call RtlAdjustPrivilege.
func RtlEnableSystemProf() (bool, error) {
	return RtlAdjustPrivilege(SESystemProf, true, false)
}

// RtlDisableSystemProf is used to disable system profile privilege that call RtlAdjustPrivilege.
func RtlDisableSystemProf() (bool, error) {
	return RtlAdjustPrivilege(SESystemProf, false, false)
}

// RtlEnableBackup is used to enable backup privilege that call RtlAdjustPrivilege.
func RtlEnableBackup() (bool, error) {
	return RtlAdjustPrivilege(SEBackup, true, false)
}

// RtlDisableBackup is used to disable backup privilege that call RtlAdjustPrivilege.
func RtlDisableBackup() (bool, error) {
	return RtlAdjustPrivilege(SEBackup, false, false)
}

// RtlEnableShutdown is used to enable shutdown privilege that call RtlAdjustPrivilege.
func RtlEnableShutdown() (bool, error) {
	return RtlAdjustPrivilege(SEShutdown, true, false)
}

// RtlDisableShutdown is used to disable shutdown privilege that call RtlAdjustPrivilege.
func RtlDisableShutdown() (bool, error) {
	return RtlAdjustPrivilege(SEShutdown, false, false)
}

// RtlEnableDebug is used to enable debug privilege that call RtlAdjustPrivilege.
func RtlEnableDebug() (bool, error) {
	return RtlAdjustPrivilege(SEDebug, true, false)
}

// RtlDisableDebug is used to disable debug privilege that call RtlAdjustPrivilege.
func RtlDisableDebug() (bool, error) {
	return RtlAdjustPrivilege(SEDebug, false, false)
}

// RtlEnableSystemEnv is used to enable system environment privilege that call RtlAdjustPrivilege.
func RtlEnableSystemEnv() (bool, error) {
	return RtlAdjustPrivilege(SESystemEnv, true, false)
}

// RtlDisableSystemEnv is used to disable system environment privilege that call RtlAdjustPrivilege.
func RtlDisableSystemEnv() (bool, error) {
	return RtlAdjustPrivilege(SESystemEnv, false, false)
}

// RtlEnableRemoteShutdown is used to enable remote shutdown privilege that call RtlAdjustPrivilege.
func RtlEnableRemoteShutdown() (bool, error) {
	return RtlAdjustPrivilege(SERemoteShutdown, true, false)
}

// RtlDisableRemoteShutdown is used to disable remote shutdown privilege that call RtlAdjustPrivilege.
func RtlDisableRemoteShutdown() (bool, error) {
	return RtlAdjustPrivilege(SERemoteShutdown, false, false)
}