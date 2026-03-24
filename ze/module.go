package ze

import (
	"runtime"

	"github.com/fumiama/gozel"
)

// ModuleHandle is a handle to a Level Zero module.
type ModuleHandle gozel.ZeModuleHandle

// ModuleCreate creates a module from SPIR-V binary data on the given device.
func (h ContextHandle) ModuleCreate(hDevice gozel.ZeDeviceHandle, data []byte) (
	ModuleHandle, error,
) {
	var m gozel.ZeModuleHandle
	_, err := gozel.ZeModuleCreate(gozel.ZeContextHandle(h), hDevice, &gozel.ZeModuleDesc{
		Stype:        gozel.ZE_STRUCTURE_TYPE_MODULE_DESC,
		Format:       gozel.ZE_MODULE_FORMAT_IL_SPIRV,
		Inputsize:    uintptr(len(data)),
		Pinputmodule: &data[0],
	}, &m, nil)
	runtime.KeepAlive(data)
	return ModuleHandle(m), err
}

// Destroy destroys the module and releases its resources.
func (h ModuleHandle) Destroy() error {
	_, err := gozel.ZeModuleDestroy(gozel.ZeModuleHandle(h))
	return err
}
