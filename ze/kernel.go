package ze

import (
	"reflect"
	"runtime"
	"unsafe"

	"github.com/fumiama/gozel"
)

// KernelHandle is a handle to a Level Zero kernel.
type KernelHandle gozel.ZeKernelHandle

// KernelCreate creates a kernel from the module by the given function name.
func (h ModuleHandle) KernelCreate(kernelName string) (KernelHandle, error) {
	b := []byte(kernelName + "\x00")
	var k gozel.ZeKernelHandle
	_, err := gozel.ZeKernelCreate(gozel.ZeModuleHandle(h), &gozel.ZeKernelDesc{
		Stype:       gozel.ZE_STRUCTURE_TYPE_KERNEL_DESC,
		Pkernelname: &b[0],
	}, &k)
	runtime.KeepAlive(b)
	return KernelHandle(k), err
}

// SetArgumentValue sets the value of a kernel argument at the given index.
func (h KernelHandle) SetArgumentValue(argIndex uint32, arg any) error {
	_, err := gozel.ZeKernelSetArgumentValue(
		gozel.ZeKernelHandle(h), argIndex, reflect.TypeOf(arg).Size(),
		*(*unsafe.Pointer)(
			unsafe.Add(unsafe.Pointer(&arg),
				unsafe.Sizeof(uintptr(0))),
		),
	)
	runtime.KeepAlive(arg)
	return err
}

// SetGroupSize sets the thread group size for the kernel.
func (h KernelHandle) SetGroupSize(groupSizeX uint32, groupSizeY uint32, groupSizeZ uint32) error {
	_, err := gozel.ZeKernelSetGroupSize(gozel.ZeKernelHandle(h), groupSizeX, groupSizeY, groupSizeZ)
	return err
}

// Destroy destroys the kernel and releases its resources.
func (h KernelHandle) Destroy() error {
	_, err := gozel.ZeKernelDestroy(gozel.ZeKernelHandle(h))
	return err
}
