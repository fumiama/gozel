package ze

import (
	"runtime"
	"unsafe"

	"github.com/fumiama/gozel/gozel"
)

// EventPoolHandle (ze_event_pool_handle_t) Handle of driver's event pool object
type EventPoolHandle gozel.ZeEventPoolHandle

// EventPoolCreate Creates a pool of events on the context.
func (h ContextHandle) EventPoolCreate(
	evcount uint32, devices ...DeviceHandle,
) (eph EventPoolHandle, err error) {
	_, err = gozel.ZeEventPoolCreate(gozel.ZeContextHandle(h), &gozel.ZeEventPoolDesc{
		Stype: gozel.ZE_STRUCTURE_TYPE_EVENT_POOL_DESC,
		Flags: gozel.ZE_EVENT_POOL_FLAG_HOST_VISIBLE,
		Count: evcount,
	}, uint32(len(devices)), (*gozel.ZeDeviceHandle)(unsafe.SliceData(devices)),
		(*gozel.ZeEventPoolHandle)(&eph),
	)
	runtime.KeepAlive(devices)
	return
}

// Destroy Deletes an event pool object.
func (h EventPoolHandle) Destroy() error {
	_, err := gozel.ZeEventPoolDestroy(gozel.ZeEventPoolHandle(h))
	return err
}

// EventHandle (ze_event_handle_t) Handle of driver's event object
type EventHandle gozel.ZeEventHandle

// EventCreate Creates an event from the pool.
func (h EventPoolHandle) EventCreate(
	index uint32, signal, wait gozel.ZeEventScopeFlags,
) (eh EventHandle, err error) {
	_, err = gozel.ZeEventCreate(gozel.ZeEventPoolHandle(h), &gozel.ZeEventDesc{
		Stype:  gozel.ZE_STRUCTURE_TYPE_EVENT_DESC,
		Index:  index,
		Signal: signal,
		Wait:   wait,
	}, (*gozel.ZeEventHandle)(&eh))
	return
}

// HostSynchronize The current host thread waits on an event to be signaled.
func (h EventHandle) HostSynchronize(timeout uint64) error {
	_, err := gozel.ZeEventHostSynchronize(gozel.ZeEventHandle(h), timeout)
	return err
}

// Destroy Deletes an event object.
func (h EventHandle) Destroy() error {
	_, err := gozel.ZeEventDestroy(gozel.ZeEventHandle(h))
	return err
}
