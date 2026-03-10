package gozel

import "unsafe"

// ZeResult maps to ze_result_t
type ZeResult uint32

const (
	ZeResultSuccess                ZeResult = 0x00000000
	ZeResultNotReady               ZeResult = 0x00000001
	ZeResultErrorUninitialized     ZeResult = 0x78000001
	ZeResultErrorInvalidArgument   ZeResult = 0x78000004
	ZeResultErrorOutOfHostMemory   ZeResult = 0x78000006
	ZeResultErrorOutOfDeviceMemory ZeResult = 0x78000007
	ZeResultErrorUnsupported       ZeResult = 0x78000009
)

func (r ZeResult) Error() string {
	switch r {
	case ZeResultSuccess:
		return "ZE_RESULT_SUCCESS"
	case ZeResultNotReady:
		return "ZE_RESULT_NOT_READY"
	case ZeResultErrorUninitialized:
		return "ZE_RESULT_ERROR_UNINITIALIZED"
	case ZeResultErrorInvalidArgument:
		return "ZE_RESULT_ERROR_INVALID_ARGUMENT"
	case ZeResultErrorOutOfHostMemory:
		return "ZE_RESULT_ERROR_OUT_OF_HOST_MEMORY"
	case ZeResultErrorOutOfDeviceMemory:
		return "ZE_RESULT_ERROR_OUT_OF_DEVICE_MEMORY"
	case ZeResultErrorUnsupported:
		return "ZE_RESULT_ERROR_UNSUPPORTED_FEATURE"
	default:
		return "ZE_RESULT_UNKNOWN"
	}
}

// ZeDriverHandle maps to ze_driver_handle_t (opaque pointer)
type ZeDriverHandle uintptr

// ZeStructureType maps to ze_structure_type_t (selected values)
type ZeStructureType uint32

const (
	// ZeStructureTypeInitDriverTypeDesc maps to ZE_STRUCTURE_TYPE_INIT_DRIVER_TYPE_DESC
	ZeStructureTypeInitDriverTypeDesc ZeStructureType = 0x00020002
)

// ZeInitDriverTypeFlags maps to ze_init_driver_type_flag_t
type ZeInitDriverTypeFlags uint32

const (
	ZeInitDriverTypeGPU ZeInitDriverTypeFlags = 1 << 0
	ZeInitDriverTypeNPU ZeInitDriverTypeFlags = 1 << 1
	// ZeInitDriverTypeAll selects all driver types (GPU + NPU)
	ZeInitDriverTypeAll ZeInitDriverTypeFlags = ZeInitDriverTypeGPU | ZeInitDriverTypeNPU
)

// ZeInitDriverTypeDesc maps to ze_init_driver_type_desc_t
type ZeInitDriverTypeDesc struct {
	Stype ZeStructureType
	PNext unsafe.Pointer
	Flags ZeInitDriverTypeFlags
}

// GPGPUDriverTypeDesc returns a ZeInitDriverTypeDesc configured for GPGPU
// (GPU driver only, correct stype).
func GPGPUDriverTypeDesc() ZeInitDriverTypeDesc {
	return ZeInitDriverTypeDesc{
		Stype: ZeStructureTypeInitDriverTypeDesc,
		Flags: ZeInitDriverTypeGPU,
	}
}
