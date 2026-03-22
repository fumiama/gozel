//go:build !linux && !windows

package zecall

// Register a process for calling. For generated init only. Not thread-safe.
func Register(string) error {
	return ErrNotImplemented
}

// Syscall invokes a registered proc by name. For generated call only.
// The go:uintptrescapes directive tells the compiler that args may contain
// pointers converted to uintptr, so the GC will keep them alive during the call.
func Syscall(name string, args ...uintptr) (uintptr, uintptr, error) {
	return 0, 0, ErrNotImplemented
}
