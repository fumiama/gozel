// Package ze provides high-level wrappers around Level Zero command objects.
package ze

import (
	"math"
	"runtime"
	"unsafe"

	"github.com/fumiama/gozel"
)

// CommandQueueHandle is a handle to a Level Zero command queue.
type CommandQueueHandle gozel.ZeCommandQueueHandle

// CommandQueueCreate creates a command queue on the given device with default mode and normal priority.
func (h ContextHandle) CommandQueueCreate(hDevice DeviceHandle) (
	CommandQueueHandle, error,
) {
	var q gozel.ZeCommandQueueHandle
	_, err := gozel.ZeCommandQueueCreate(gozel.ZeContextHandle(h), gozel.ZeDeviceHandle(hDevice), &gozel.ZeCommandQueueDesc{
		Stype:    gozel.ZE_STRUCTURE_TYPE_COMMAND_QUEUE_DESC,
		Mode:     gozel.ZE_COMMAND_QUEUE_MODE_DEFAULT,
		Priority: gozel.ZE_COMMAND_QUEUE_PRIORITY_NORMAL,
	}, &q)
	return CommandQueueHandle(q), err
}

// ExecuteCommandLists submits the command list for execution on the command queue.
func (h CommandQueueHandle) ExecuteCommandLists(hCommandList ...CommandListHandle) error {
	_, err := gozel.ZeCommandQueueExecuteCommandLists(
		gozel.ZeCommandQueueHandle(h), uint32(len(hCommandList)),
		(*gozel.ZeCommandListHandle)(&hCommandList[0]), 0,
	)
	runtime.KeepAlive(hCommandList)
	return err
}

// Synchronize blocks the host until all commands in the command queue have completed.
func (h CommandQueueHandle) Synchronize() error {
	_, err := gozel.ZeCommandQueueSynchronize(gozel.ZeCommandQueueHandle(h), math.MaxUint64)
	return err
}

// Destroy destroys the command queue and releases its resources.
func (h CommandQueueHandle) Destroy() error {
	_, err := gozel.ZeCommandQueueDestroy(gozel.ZeCommandQueueHandle(h))
	return err
}

// CommandListHandle is a handle to a Level Zero command list.
type CommandListHandle gozel.ZeCommandListHandle

// CommandListCreate creates a command list on the given device.
func (h ContextHandle) CommandListCreate(hDevice DeviceHandle) (
	CommandListHandle, error,
) {
	var cl gozel.ZeCommandListHandle
	_, err := gozel.ZeCommandListCreate(gozel.ZeContextHandle(h), gozel.ZeDeviceHandle(hDevice), &gozel.ZeCommandListDesc{
		Stype: gozel.ZE_STRUCTURE_TYPE_COMMAND_LIST_DESC,
	}, &cl)
	return CommandListHandle(cl), err
}

// AppendLaunchKernel appends a kernel launch command to the command list.
func (h CommandListHandle) AppendLaunchKernel(
	hKernel KernelHandle, pLaunchFuncArgs *gozel.ZeGroupCount,
) error {
	_, err := gozel.ZeCommandListAppendLaunchKernel(gozel.ZeCommandListHandle(h), gozel.ZeKernelHandle(hKernel), pLaunchFuncArgs, 0, 0, nil)
	return err
}

// AppendLaunchKernel appends a kernel launch command to the command list with args.
func (h CommandListHandle) AppendLaunchKernelWithArguments(
	hCommandList gozel.ZeCommandListHandle, hKernel gozel.ZeKernelHandle, groupCounts *gozel.ZeGroupCount,
	groupSizes *gozel.ZeGroupSize, pArguments *unsafe.Pointer,
) error {
	_, err := gozel.ZeCommandListAppendLaunchKernelWithArguments(
		gozel.ZeCommandListHandle(h), gozel.ZeKernelHandle(hKernel), groupCounts, groupSizes, pArguments, nil, 0, 0, nil)
	return err
}

// Close closes the command list, making it ready for execution.
func (h CommandListHandle) Close() error {
	_, err := gozel.ZeCommandListClose(gozel.ZeCommandListHandle(h))
	return err
}

// AppendMemoryCopy appends a memory copy command from srcptr to dstptr of the given size.
func (h CommandListHandle) AppendMemoryCopy(
	dstptr unsafe.Pointer, srcptr unsafe.Pointer, size uintptr,
) error {
	_, err := gozel.ZeCommandListAppendMemoryCopy(gozel.ZeCommandListHandle(h), dstptr, srcptr, size, 0, 0, nil)
	return err
}

// Destroy destroys the command list and releases its resources.
func (h CommandListHandle) Destroy() error {
	_, err := gozel.ZeCommandListDestroy(gozel.ZeCommandListHandle(h))
	return err
}

// AppendBarrier appends an execution barrier to the command list.
func (h CommandListHandle) AppendBarrier() error {
	_, err := gozel.ZeCommandListAppendBarrier(gozel.ZeCommandListHandle(h), 0, 0, nil)
	return err
}
