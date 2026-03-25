package ze

import "github.com/fumiama/gozel"

// DeviceHandle is a handle to a Level Zero driver's device object.
type DeviceHandle gozel.ZeDeviceHandle

// DeviceGet retrieves all devices within the driver.
func (h DriverHandle) DeviceGet() ([]DeviceHandle, error) {
	var count uint32
	_, err := gozel.ZeDeviceGet(gozel.ZeDriverHandle(h), &count, nil)
	if err != nil {
		return nil, err
	}
	if count == 0 {
		return nil, nil
	}
	devices := make([]DeviceHandle, count)
	_, err = gozel.ZeDeviceGet(gozel.ZeDriverHandle(h), &count, (*gozel.ZeDeviceHandle)(&devices[0]))
	if err != nil {
		return nil, err
	}
	return devices, nil
}

// DeviceGetProperties retrieves properties of the device.
func (h DeviceHandle) DeviceGetProperties() (prop gozel.ZeDeviceProperties, err error) {
	prop.Stype = gozel.ZE_STRUCTURE_TYPE_DEVICE_PROPERTIES
	_, err = gozel.ZeDeviceGetProperties(gozel.ZeDeviceHandle(h), &prop)
	return
}

// DeviceGetComputeProperties retrieves compute properties of the device.
func (h DeviceHandle) DeviceGetComputeProperties() (prop gozel.ZeDeviceComputeProperties, err error) {
	prop.Stype = gozel.ZE_STRUCTURE_TYPE_DEVICE_COMPUTE_PROPERTIES
	_, err = gozel.ZeDeviceGetComputeProperties(gozel.ZeDeviceHandle(h), &prop)
	return
}
