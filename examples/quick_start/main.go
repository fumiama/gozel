// Package main demonstrates quick start usage of the gozel Level Zero bindings.
package main

import (
	"fmt"
	"strings"

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
		name, _, _ := strings.Cut(string(prop.Name[:]), "\x00")
		fmt.Printf("  Device: %s\n", name)
	}

	// Found 1 GPU driver(s)
	//   Device:  Graphics
}
