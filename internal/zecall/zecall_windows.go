package gozel

import (
	"errors"
	"syscall"
)

const (
	zeLibraryName = "ze_loader.dll"
)

var (
	// ErrZeCallNotInit please call Init() first
	ErrZeCallNotInit = errors.New("zecall not init")
	// ErrNoSuchProcess please register the process first
	ErrNoSuchProcess = errors.New("no such process")
)

var (
	libZeLoader *syscall.DLL
	procMap     = map[string]*syscall.Proc{}
)

// Init load lib using syscall
func Init() error {
	h, err := syscall.LoadLibrary(zeLibraryName)
	if err != nil {
		return err
	}
	libZeLoader = &syscall.DLL{Handle: h, Name: zeLibraryName}

	return nil
}

// Register a process for calling
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

// Call a process
func Call(name string, args ...uintptr) (r1, r2 uintptr, err error) {
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
