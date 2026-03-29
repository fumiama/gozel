package ze

import "github.com/fumiama/gozel/gozel"

// SamplerHandle (ze_sampler_handle_t) Handle of driver's sampler object
type SamplerHandle gozel.ZeSamplerHandle

// SamplerCreate Creates sampler on the context.
func (h ContextHandle) SamplerCreate(
	hDevice DeviceHandle, addressmode gozel.ZeSamplerAddressMode,
	filtermode gozel.ZeSamplerFilterMode, isnormalized gozel.ZeBool,
) (sh SamplerHandle, err error) {
	_, err = gozel.ZeSamplerCreate(gozel.ZeContextHandle(h), gozel.ZeDeviceHandle(hDevice), &gozel.ZeSamplerDesc{
		Stype:        gozel.ZE_STRUCTURE_TYPE_SAMPLER_DESC,
		Addressmode:  addressmode,
		Filtermode:   filtermode,
		Isnormalized: isnormalized,
	}, (*gozel.ZeSamplerHandle)(&sh))
	return
}

// Destroy Destroys sampler object.
func (h SamplerHandle) Destroy() error {
	_, err := gozel.ZeSamplerDestroy(gozel.ZeSamplerHandle(h))
	return err
}
