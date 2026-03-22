package zecall

import (
	"errors"
	"math"
	"reflect"
)

type ReturnTypes interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 |
		~uint16 | ~uint32 | ~uint64 | ~uintptr | ~float32 | ~float64
}

// Call invokes a registered proc by name. For generated call only.
// The go:uintptrescapes directive tells the compiler that args may contain
// pointers converted to uintptr, so the GC will keep them alive during the call.
//
//go:uintptrescapes
func Call[T ReturnTypes](name string, args ...uintptr) (r T, err error) {
	r1, r2, err := Syscall(name, args...)
	if err != nil {
		return
	}
	k := reflect.TypeOf(r).Kind()
	switch k {
	case reflect.Int, reflect.Int8, reflect.Int16,
		reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8,
		reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		r = (T)(r1)
	case reflect.Float32:
		r = (T)(math.Float32frombits(uint32(r2)))
	case reflect.Float64:
		r = (T)(math.Float64frombits(uint64(r2)))
	default:
		err = errors.New("zecall unsupported kind " + k.String())
	}
	return
}
