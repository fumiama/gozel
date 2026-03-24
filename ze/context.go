package ze

import "github.com/fumiama/gozel"

// ContextHandle is a handle to a Level Zero context.
type ContextHandle gozel.ZeContextHandle

// ContextCreate creates a new context for the driver.
func (h DriverHandle) ContextCreate() (ContextHandle, error) {
	var ctx gozel.ZeContextHandle
	_, err := gozel.ZeContextCreate(gozel.ZeDriverHandle(h), &gozel.ZeContextDesc{
		Stype: gozel.ZE_STRUCTURE_TYPE_CONTEXT_DESC,
	}, &ctx)
	return ContextHandle(ctx), err
}

// Destroy destroys the context and releases its resources.
func (h ContextHandle) Destroy() error {
	_, err := gozel.ZeContextDestroy(gozel.ZeContextHandle(h))
	return err
}
