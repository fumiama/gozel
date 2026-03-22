package zecall

import (
	"syscall"

	"github.com/ebitengine/purego"
)

const (
	zeLibraryName = "libze_loader.so"
)

var (
	libZeLoader uintptr
	procMap     = map[string]uintptr{}
)

func init() {
	if libZeLoader != 0 {
		return
	}
	h, err := purego.Dlopen(zeLibraryName, purego.RTLD_LAZY|purego.RTLD_GLOBAL)
	if err != nil {
		panic(err)
	}
	libZeLoader = h
}

// Register a process for calling. For generated init only. Not thread-safe.
func Register(name string) error {
	if libZeLoader == 0 {
		return ErrZeCallNotInit
	}
	sym, err := purego.Dlsym(libZeLoader, name)
	if err != nil {
		return err
	}
	procMap[name] = sym
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
	var errno uintptr
	r1, r2, errno = purego.SyscallN(fn, args...)
	if r1 == 0 {
		return
	}
	if errno != 0 {
		err = syscall.Errno(errno)
	}
	return
}
