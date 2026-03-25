// Package main demonstrates quick start usage of the gozel Level Zero bindings.
package main

import (
	"fmt"

	"github.com/fumiama/gozel/ze"
)

func main() {
	gpus, err := ze.InitGPUDrivers()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Found %d GPU driver(s)\n", len(gpus))

	devs, _ := gpus[0].DeviceGet()
	for _, d := range devs {
		prop, _ := d.DeviceGetProperties()
		fmt.Printf("  Device: %s\n", string(prop.Name[:]))
	}

	// Found 1 GPU driver(s)
	//   Device:  Graphics
}
