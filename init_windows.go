package gozel

import (
	"fmt"
	"syscall"
)

const (
	zeLibraryName = "ze_loader.dll"
)

var (
	libZeLoader syscall.DLL
)

func InitZe() error {
	h, err := syscall.LoadLibrary(zeLibraryName)
	if err != nil {
		return err
	}
	libZeLoader = syscall.DLL{Handle: h, Name: zeLibraryName}
	procZeInitDrivers, err = libZeLoader.FindProc("zeInitDrivers")
	if err != nil {
		return fmt.Errorf("zeInitDrivers not found in ze_loader.dll: %w", err)
	}
	return nil
}
