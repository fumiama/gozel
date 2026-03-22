package zecall

import (
	"syscall"
)

const (
	zeLibraryName = "ze_loader.dll"
)

var (
	libZeLoader *syscall.DLL
	procMap     = map[string]*syscall.Proc{}
)

func init() {
	if libZeLoader != nil {
		return
	}
	h, err := syscall.LoadLibrary(zeLibraryName)
	if err != nil {
		panic(err)
	}
	libZeLoader = &syscall.DLL{Handle: h, Name: zeLibraryName}
}

// Register a process for calling. For generated init only. Not thread-safe.
func Register(name string) error {
	if libZeLoader == nil {
		return ErrZeCallNotInit
	}
	proc, err := libZeLoader.FindProc(name)
	if err != nil {
		return err
	}
	procMap[name] = proc
	return nil
}

// Syscall invokes a registered proc by name. For generated call only.
// The go:uintptrescapes directive tells the compiler that args may contain
// pointers converted to uintptr, so the GC will keep them alive during the call.
//
//go:uintptrescapes
func Syscall(name string, args ...uintptr) (r1, r2 uintptr, err error) {
	fn, ok := procMap[name]
	if !ok {
		return 0, 0, ErrNoSuchProcess
	}
	r1, r2, err = fn.Call(args...)
	if r1 == 0 {
		err = nil
	}
	return
}
