package ze

import (
	"errors"
	"runtime"
	"strings"

	"github.com/fumiama/gozel"
)

// ModuleHandle is a handle to a Level Zero module.
type ModuleHandle gozel.ZeModuleHandle

// ModuleCreate creates a module from SPIR-V binary data on the given device.
func (h ContextHandle) ModuleCreate(hDevice gozel.ZeDeviceHandle, data []byte) (
	ModuleHandle, error,
) {
	var (
		m  gozel.ZeModuleHandle
		lg gozel.ZeModuleBuildLogHandle
	)
	_, err := gozel.ZeModuleCreate(gozel.ZeContextHandle(h), hDevice, &gozel.ZeModuleDesc{
		Stype:        gozel.ZE_STRUCTURE_TYPE_MODULE_DESC,
		Format:       gozel.ZE_MODULE_FORMAT_IL_SPIRV,
		Inputsize:    uintptr(len(data)),
		Pinputmodule: &data[0],
	}, &m, &lg)
	runtime.KeepAlive(data)
	defer gozel.ZeModuleBuildLogDestroy(lg)
	if err != nil {
		var lgsz uintptr
		_, errlg := gozel.ZeModuleBuildLogGetString(lg, &lgsz, nil)
		if errlg == nil {
			data := make([]byte, lgsz)
			_, errlg := gozel.ZeModuleBuildLogGetString(lg, &lgsz, &data[0])
			runtime.KeepAlive(data)
			if errlg == nil {
				sb := strings.Builder{}
				sb.WriteString(err.Error())
				sb.WriteString(", build log: ")
				sb.Write(data)
				err = errors.New(sb.String())
			}
		}
	}
	return ModuleHandle(m), err
}

// Destroy destroys the module and releases its resources.
func (h ModuleHandle) Destroy() error {
	_, err := gozel.ZeModuleDestroy(gozel.ZeModuleHandle(h))
	return err
}
