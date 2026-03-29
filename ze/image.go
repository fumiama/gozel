package ze

import "github.com/fumiama/gozel/gozel"

type ImageHandle gozel.ZeImageHandle

// ImageCreate Creates a 2D image on the context.
// flags: 0 for read-only (kernel input), ZE_IMAGE_FLAG_KERNEL_WRITE for writable (kernel output).
func (h ContextHandle) ImageCreate(
	hDevice DeviceHandle, flags gozel.ZeImageFlags, format gozel.ZeImageFormat,
	width uint64, height uint32,
) (ih ImageHandle, err error) {
	_, err = gozel.ZeImageCreate(gozel.ZeContextHandle(h), gozel.ZeDeviceHandle(hDevice),
		&gozel.ZeImageDesc{
			Stype:  gozel.ZE_STRUCTURE_TYPE_IMAGE_DESC,
			Flags:  flags,
			Type:   gozel.ZE_IMAGE_TYPE_2D,
			Format: format,
			Width:  width,
			Height: height,
		}, (*gozel.ZeImageHandle)(&ih))
	return
}

// Destroy Deletes an image object.
func (h ImageHandle) Destroy() error {
	_, err := gozel.ZeImageDestroy(gozel.ZeImageHandle(h))
	return err
}
