package main

import (
	_ "embed"
	"fmt"
	"math/rand"
	"os"
	"unsafe"

	"github.com/fumiama/gozel"
	"github.com/fumiama/gozel/ze"
)

//go:generate clang++ -fsycl -fsycl-device-only -fno-sycl-use-footer -faddrsig -Xclang -emit-llvm-bc main.cpp -o device_func.bc
//go:generate sycl-post-link -symbols -split=auto -o device_func.table device_func.bc
//go:generate llvm-spirv -o device_func.spv device_func_0.bc
//go:generate clang++ -target spir64-unknown-unknown -S -emit-llvm -x ir device_func_0.bc -o device_func.ll
//go:generate go run ../../cmd/func2kernel device_func.ll device_kern.ll
//go:generate clang++ -target spir64-unknown-unknown -c -emit-llvm -x ir device_kern.ll -o device_kern.bc
//go:generate llvm-spirv -o main.spv device_kern.bc
//go:generate clang++ -target spir64-unknown-unknown -S -emit-llvm -x ir device_kern.bc -o main.ll

//go:embed main.spv
var kernelspv []byte

const (
	X, Y, Z = 1024, 1, 1
	N       = X * Y * Z
	bufsz   = N * unsafe.Sizeof(float64(0))
)

func main() {
	floatbuf := make([]float64, 2*N)
	for i := range floatbuf {
		floatbuf[i] = rand.Float64()
	}

	gpus, err := ze.InitGPUDrivers()
	if err != nil {
		panic(err)
	}
	if len(gpus) == 0 {
		panic("no gpu available")
	}
	gpu := gpus[0]

	ctx, err := gpu.ContextCreate()
	if err != nil {
		panic(err)
	}

	devs, err := gpu.DeviceGet()
	if err != nil {
		panic(err)
	}
	if len(devs) == 0 {
		panic("no device available")
	}
	dev := devs[0]

	q, err := ctx.CommandQueueCreate(dev)
	if err != nil {
		panic(err)
	}
	defer q.Destroy()

	hbuf_v1, err := ctx.MemAllocHost(bufsz, 1)
	if err != nil {
		panic(err)
	}
	defer ctx.MemFree(hbuf_v1)

	hbuf_v2, err := ctx.MemAllocHost(bufsz, 1)
	if err != nil {
		panic(err)
	}
	defer ctx.MemFree(hbuf_v2)

	dbuf_v1, err := ctx.MemAllocDevice(dev, bufsz, 1)
	if err != nil {
		panic(err)
	}
	defer ctx.MemFree(dbuf_v1)

	dbuf_v2, err := ctx.MemAllocDevice(dev, bufsz, 1)
	if err != nil {
		panic(err)
	}
	defer ctx.MemFree(dbuf_v2)

	zev1, zev2 := unsafe.Slice((*float64)(hbuf_v1), N), unsafe.Slice((*float64)(hbuf_v2), N)
	copy(zev1, floatbuf[:N])
	copy(zev2, floatbuf[N:])

	mod, err := ctx.ModuleCreate(dev, kernelspv)
	if err != nil {
		panic(err)
	}
	defer mod.Destroy()

	krn, err := mod.KernelCreate("vector_add")
	if err != nil {
		panic(err)
	}
	defer krn.Destroy()

	err = krn.SetArgumentValue(0, unsafe.Sizeof(uintptr(0)), unsafe.Pointer(&dbuf_v1))
	if err != nil {
		panic(err)
	}
	err = krn.SetArgumentValue(1, unsafe.Sizeof(uintptr(0)), unsafe.Pointer(&dbuf_v2))
	if err != nil {
		panic(err)
	}
	err = krn.SetGroupSize(X, Y, Z)
	if err != nil {
		panic(err)
	}

	lst, err := ctx.CommandListCreate(dev)
	if err != nil {
		panic(err)
	}
	defer lst.Destroy()

	err = lst.AppendMemoryCopy(dbuf_v1, hbuf_v1, bufsz)
	if err != nil {
		panic(err)
	}
	err = lst.AppendMemoryCopy(dbuf_v2, hbuf_v2, bufsz)
	if err != nil {
		panic(err)
	}

	err = lst.AppendBarrier()
	if err != nil {
		panic(err)
	}

	err = lst.AppendLaunchKernel(krn, &gozel.ZeGroupCount{
		Groupcountx: 1, Groupcounty: 1, Groupcountz: 1,
	})
	if err != nil {
		panic(err)
	}

	err = lst.AppendBarrier()
	if err != nil {
		panic(err)
	}

	err = lst.AppendMemoryCopy(hbuf_v1, dbuf_v1, bufsz)
	if err != nil {
		panic(err)
	}

	err = lst.Close()
	if err != nil {
		panic(err)
	}

	err = q.ExecuteCommandLists(lst)
	if err != nil {
		panic(err)
	}

	err = q.Synchronize()
	if err != nil {
		panic(err)
	}

	fail := false
	for i := range N {
		expect := floatbuf[i] + floatbuf[N+i]
		if zev1[i] != expect {
			fail = true
			fmt.Printf("[%05d] expect %f = %f + %f, got %f.\n", i, expect, floatbuf[i], floatbuf[N+i], zev1[i])
		} else {
			fmt.Printf("[%05d] valid  %f = %f + %f, got %f.\n", i, expect, floatbuf[i], floatbuf[N+i], zev1[i])
		}
	}

	if fail {
		os.Exit(1)
	}
}
