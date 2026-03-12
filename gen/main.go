package main

import (
	"archive/zip"
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os"
	"strings"
)

func main() {
	spec := flag.String("spec", "v1.28.2", "The l0 loader spec version tag starting with v or a local level-zero path for dev.")
	flag.Parse()

	var specdir fs.FS

	if strings.HasPrefix(*spec, "v") {
		ver := (*spec)[1:]
		u := fmt.Sprintf("https://github.com/oneapi-src/level-zero/releases/download/v%s/level-zero-win-sdk-%s.zip", ver, ver)
		resp, err := http.Get(u)
		if err != nil {
			panic(err)
		}
		data, err := io.ReadAll(resp.Body)
		_ = resp.Body.Close()
		if err != nil {
			panic(err)
		}
		r, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
		if err != nil {
			panic(err)
		}
		specdir = r
	} else {
		specdir = os.DirFS(*spec)
	}

	fmt.Println("[main] parsing core APIs...")
	f, err := specdir.Open("include/level_zero/ze_api.h")
	if err != nil {
		panic(err)
	}
	scanheader("core", bufio.NewScanner(f))
	_ = f.Close()
	fmt.Println("[main] finish parsing core")

	fmt.Println("[main] parsing runtime APIs...")
	f, err = specdir.Open("include/level_zero/zer_api.h")
	if err != nil {
		panic(err)
	}
	scanheader("rntm", bufio.NewScanner(f))
	_ = f.Close()
	fmt.Println("[main] finish parsing runtime")

	fmt.Println("[main] parsing sysman APIs...")
	f, err = specdir.Open("include/level_zero/zes_api.h")
	if err != nil {
		panic(err)
	}
	scanheader("sysm", bufio.NewScanner(f))
	_ = f.Close()
	fmt.Println("[main] finish parsing sysman")

	fmt.Println("[main] parsing tools APIs...")
	f, err = specdir.Open("include/level_zero/zet_api.h")
	if err != nil {
		panic(err)
	}
	scanheader("tols", bufio.NewScanner(f))
	_ = f.Close()
	fmt.Println("[main] finish parsing tools")
}
