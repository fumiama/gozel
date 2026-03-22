package zecall

import (
	"errors"
	"runtime"
)

var (
	// ErrNotImplemented is a stub error.
	ErrNotImplemented = errors.New("zecall is not implemtent on" + runtime.GOOS + " " + runtime.GOARCH)
	// ErrZeCallNotInit please call Init() first.
	ErrZeCallNotInit = errors.New("zecall not init")
	// ErrNoSuchProcess please register the process first.
	ErrNoSuchProcess = errors.New("no such process")
)
