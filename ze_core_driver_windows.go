package gozel

import (
	"syscall"
	"unsafe"
)

var procZeInitDrivers *syscall.Proc

// InitDrivers calls zeInitDrivers from ze_loader.dll.
// On success pCount contains the number of drivers and phDrivers (if non-nil)
// is filled with driver handles.
func InitDrivers(desc *ZeInitDriverTypeDesc) ([]ZeDriverHandle, error) {
	var count uint32
	r, _, _ := procZeInitDrivers.Call(
		uintptr(unsafe.Pointer(&count)),
		0,
		uintptr(unsafe.Pointer(desc)),
	)
	if ZeResult(r) != ZeResultSuccess {
		return nil, ZeResult(r)
	}
	if count == 0 {
		return nil, nil
	}

	handles := make([]ZeDriverHandle, count)
	r, _, _ = procZeInitDrivers.Call(
		uintptr(unsafe.Pointer(&count)),
		uintptr(unsafe.Pointer(&handles[0])),
		uintptr(unsafe.Pointer(desc)),
	)
	if ZeResult(r) != ZeResultSuccess {
		return nil, ZeResult(r)
	}
	return handles, nil
}
