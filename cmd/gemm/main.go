package main

import (
	"fmt"

	"github.com/fumiama/gozel"
)

func main() {
	err := gozel.InitZe()
	if err != nil {
		panic(err)
	}
	desc := gozel.GPGPUDriverTypeDesc()
	fmt.Println(gozel.InitDrivers(&desc))
}
