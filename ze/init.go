package ze

import (
	"github.com/fumiama/gozel"
)

func initDrivers(flags gozel.ZeInitDriverTypeFlags) ([]gozel.ZeDriverHandle, error) {
	var count uint32
	desc := &gozel.ZeInitDriverTypeDesc{
		Stype: gozel.ZE_STRUCTURE_TYPE_INIT_DRIVER_TYPE_DESC,
		Flags: flags,
	}
	_, err := gozel.ZeInitDrivers(&count, nil, desc)
	if count == 0 {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	handles := make([]gozel.ZeDriverHandle, count)
	_, err = gozel.ZeInitDrivers(&count, &handles[0], desc)
	if err != nil {
		return nil, err
	}
	return handles, nil
}

// InitGPUDrivers calls zeInitDrivers with ZE_INIT_DRIVER_TYPE_FLAG_GPU from ze_loader.dll.
// On success pCount contains the number of drivers and phDrivers (if non-nil)
// is filled with driver handles.
func InitGPUDrivers() ([]gozel.ZeDriverHandle, error) {
	return initDrivers(gozel.ZE_INIT_DRIVER_TYPE_FLAG_GPU)
}

// InitNPUDrivers calls zeInitDrivers with ZE_INIT_DRIVER_TYPE_FLAG_NPU from ze_loader.dll.
// On success pCount contains the number of drivers and phDrivers (if non-nil)
// is filled with driver handles.
func InitNPUDrivers() ([]gozel.ZeDriverHandle, error) {
	return initDrivers(gozel.ZE_INIT_DRIVER_TYPE_FLAG_NPU)
}
