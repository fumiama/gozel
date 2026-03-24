package ze

import "github.com/fumiama/gozel"

// DeviceGet retrieves all devices within the driver.
func (h DriverHandle) DeviceGet() ([]gozel.ZeDeviceHandle, error) {
	var count uint32
	_, err := gozel.ZeDeviceGet(gozel.ZeDriverHandle(h), &count, nil)
	if err != nil {
		return nil, err
	}
	if count == 0 {
		return nil, nil
	}
	devices := make([]gozel.ZeDeviceHandle, count)
	_, err = gozel.ZeDeviceGet(gozel.ZeDriverHandle(h), &count, &devices[0])
	if err != nil {
		return nil, err
	}
	return devices, nil
}
