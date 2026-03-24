package main

import (
	"bufio"
	"os"
	"strings"
)

func main() {
	f, err := os.Open(os.Args[1])
	if err != nil {
		panic(err)
	}
	defer f.Close()
	fo, err := os.Create(os.Args[2])
	if err != nil {
		panic(err)
	}
	defer fo.Close()
	scan := bufio.NewScanner(f)
	for scan.Scan() {
		t := scan.Text()
		t = strings.ReplaceAll(t, " spir_func ", " spir_kernel ")
		t = strings.ReplaceAll(t, "ptr addrspace(4)", "ptr addrspace(1)")
		fo.WriteString(t)
		fo.WriteString("\n")
	}
}
