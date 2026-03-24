package ze

import (
	"unsafe"

	"github.com/fumiama/gozel"
)

// MemAllocDevice allocates device memory on the given device with the specified size and alignment.
func (h ContextHandle) MemAllocDevice(hDevice gozel.ZeDeviceHandle, size uintptr, alignment uintptr) (
	unsafe.Pointer, error,
) {
	var p unsafe.Pointer
	_, err := gozel.ZeMemAllocDevice(gozel.ZeContextHandle(h), &gozel.ZeDeviceMemAllocDesc{
		Stype: gozel.ZE_STRUCTURE_TYPE_DEVICE_MEM_ALLOC_DESC,
	}, size, alignment, hDevice, &p)
	return p, err
}

// MemAllocHost allocates host memory with the specified size and alignment.
func (h ContextHandle) MemAllocHost(size uintptr, alignment uintptr) (
	unsafe.Pointer, error,
) {
	var p unsafe.Pointer
	_, err := gozel.ZeMemAllocHost(gozel.ZeContextHandle(h), &gozel.ZeHostMemAllocDesc{
		Stype: gozel.ZE_STRUCTURE_TYPE_DEVICE_MEM_ALLOC_DESC,
	}, size, alignment, &p)
	return p, err
}

// MemFree frees memory previously allocated with MemAllocDevice or MemAllocHost.
func (h ContextHandle) MemFree(ptr unsafe.Pointer) error {
	_, err := gozel.ZeMemFree(gozel.ZeContextHandle(h), ptr)
	return err
}
