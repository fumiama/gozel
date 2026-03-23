package main

import (
	"fmt"

	"github.com/fumiama/gozel/ze"
)

func main() {
	hs, err := ze.InitGPUDrivers()
	if err != nil {
		panic(err)
	}
	fmt.Println(hs)
}
