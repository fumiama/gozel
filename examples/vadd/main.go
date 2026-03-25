// Package main demonstrates vector addition using the gozel Level Zero bindings.
package main

import (
	_ "embed"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
	"unsafe"

	"github.com/fumiama/gozel"
	"github.com/fumiama/gozel/ze"
)

//go:generate clang++ -fsycl -fsycl-device-only -fsycl-targets=spirv64 -faddrsig -Xclang -emit-llvm-bc main.cpp -o device_func.bc
//go:generate sycl-post-link -symbols -split=auto -o device_func.table device_func.bc
//go:generate clang++ -target spirv64-unknown-unknown -S -emit-llvm -x ir device_func_0.bc -o device_func.ll
//go:generate go run ../../cmd/func2kernel device_func.ll device_kern.ll
//go:generate clang++ -target spirv64-unknown-unknown -c -emit-llvm -x ir device_kern.ll -o device_kern.bc
//go:generate llvm-spirv -o main.spv device_kern.bc
//go:generate clang++ -target spirv64-unknown-unknown -S -emit-llvm -x ir device_kern.bc -o main.ll

//go:embed main.spv
var kernelspv []byte

func main() {
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

	prop, err := dev.DeviceGetProperties()
	if err != nil {
		panic(err)
	}

	fmt.Println("===============  Device Basic Properties  ===============")
	fmt.Println(
		"Running on device: ID =", prop.Deviceid, ", Name =",
		strings.TrimSpace(string(prop.Name[:])),
		"@", strconv.FormatFloat(float64(prop.Coreclockrate)/1024/1024/1024, 'f', 2, 64), "GHz.",
	)

	cprop, err := dev.DeviceGetComputeProperties()
	if err != nil {
		panic(err)
	}
	fmt.Println("=============== Device Compute Properties ===============")
	fmt.Printf("%-28s (%d, %d, %d)\n", "Max Group Size (X, Y, Z):", cprop.Maxgroupsizex, cprop.Maxgroupsizey, cprop.Maxgroupsizez)
	fmt.Printf("%-28s (%d, %d, %d)\n", "Max Group Count (X, Y, Z):", cprop.Maxgroupcountx, cprop.Maxgroupcounty, cprop.Maxgroupcountz)
	fmt.Printf("%-28s %d\n", "Max Total Group Size:", cprop.Maxtotalgroupsize)
	fmt.Printf("%-28s %d\n", "Max Shared Local Memory:", cprop.Maxsharedlocalmemory)
	fmt.Printf("%-28s %d\n", "Num Subgroup Sizes:", cprop.Numsubgroupsizes)
	fmt.Printf("%-28s %v\n", "Subgroup Sizes:", cprop.Subgroupsizes[:])

	var (
		X, Y, Z    = uintptr(cprop.Maxgroupsizex), uintptr(1), uintptr(1)
		groupCount = uintptr(65536)
		N          = X * groupCount
		bufsz      = N * unsafe.Sizeof(float32(0))
	)
	fmt.Println("=============== Computation Configuration ===============")
	fmt.Printf("%-28s (%d, %d, %d)\n", "Group Size (X, Y, Z):", X, Y, Z)
	fmt.Printf("%-28s %d\n", "Group Count:", groupCount)
	fmt.Printf("%-28s %d\n", "Total Elements (N):", N)
	fmt.Printf("%-28s %d MiB\n", "Buffer Size:", bufsz/1024/1024)

	q, err := ctx.CommandQueueCreate(dev)
	if err != nil {
		panic(err)
	}
	defer q.Destroy()

	hbufV1, err := ctx.MemAllocHost(bufsz, 1)
	if err != nil {
		panic(err)
	}
	defer ctx.MemFree(hbufV1)

	hbufV2, err := ctx.MemAllocHost(bufsz, 1)
	if err != nil {
		panic(err)
	}
	defer ctx.MemFree(hbufV2)

	dbufV1, err := ctx.MemAllocDevice(dev, bufsz, 1)
	if err != nil {
		panic(err)
	}
	defer ctx.MemFree(dbufV1)

	dbufV2, err := ctx.MemAllocDevice(dev, bufsz, 1)
	if err != nil {
		panic(err)
	}
	defer ctx.MemFree(dbufV2)

	floatbuf := make([]float32, 2*N)
	for i := range floatbuf {
		floatbuf[i] = rand.Float32()
	}

	zev1, zev2 := unsafe.Slice((*float32)(hbufV1), N), unsafe.Slice((*float32)(hbufV2), N)
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

	err = krn.SetArgumentValue(0, unsafe.Sizeof(uintptr(0)), unsafe.Pointer(&dbufV1))
	if err != nil {
		panic(err)
	}
	err = krn.SetArgumentValue(1, unsafe.Sizeof(uintptr(0)), unsafe.Pointer(&dbufV2))
	if err != nil {
		panic(err)
	}
	err = krn.SetGroupSize(uint32(X), uint32(Y), uint32(Z))
	if err != nil {
		panic(err)
	}

	lstpre, err := ctx.CommandListCreate(dev)
	if err != nil {
		panic(err)
	}
	defer lstpre.Destroy()

	err = lstpre.AppendMemoryCopy(dbufV1, hbufV1, bufsz)
	if err != nil {
		panic(err)
	}
	err = lstpre.AppendMemoryCopy(dbufV2, hbufV2, bufsz)
	if err != nil {
		panic(err)
	}

	err = lstpre.AppendBarrier()
	if err != nil {
		panic(err)
	}

	err = lstpre.Close()
	if err != nil {
		panic(err)
	}

	lstcalc, err := ctx.CommandListCreate(dev)
	if err != nil {
		panic(err)
	}
	defer lstcalc.Destroy()

	err = lstcalc.AppendLaunchKernel(krn, &gozel.ZeGroupCount{
		Groupcountx: uint32(groupCount), Groupcounty: 1, Groupcountz: 1,
	})
	if err != nil {
		panic(err)
	}

	err = lstcalc.AppendBarrier()
	if err != nil {
		panic(err)
	}

	err = lstcalc.Close()
	if err != nil {
		panic(err)
	}

	lstpost, err := ctx.CommandListCreate(dev)
	if err != nil {
		panic(err)
	}
	defer lstpost.Destroy()

	err = lstpost.AppendMemoryCopy(hbufV1, dbufV1, bufsz)
	if err != nil {
		panic(err)
	}

	err = lstpost.Close()
	if err != nil {
		panic(err)
	}

	start := time.Now()
	err = q.ExecuteCommandLists(lstpre, lstcalc, lstpost)
	if err != nil {
		panic(err)
	}
	err = q.Synchronize()
	if err != nil {
		panic(err)
	}
	elapsed := time.Since(start)

	fmt.Println("===============    Calculation Results    ===============")
	fmt.Printf("%-28s %.6f ms\n", "GPU Execution Time:", elapsed.Seconds()*1000)
	fmt.Printf("%-28s %.2f GiB/s\n", "GPU Throughput:", float64(bufsz)/elapsed.Seconds()/1e9)

	tmpbuf := make([]float32, N)
	start = time.Now()
	for i := range N {
		tmpbuf[i] = floatbuf[i] + floatbuf[N+i]
	}
	elapsed = time.Since(start)

	fmt.Println("===============    Validation Results    ===============")
	fmt.Printf("%-28s %.6f ms\n", "CPU Execution Time:", elapsed.Seconds()*1000)
	fmt.Printf("%-28s %.2f GiB/s\n", "CPU Throughput:", float64(bufsz)/elapsed.Seconds()/1e9)

	fail := false
	for i := range N {
		expect := floatbuf[i] + floatbuf[N+i]
		if zev1[i] != expect {
			fail = true
			fmt.Printf("[%05d] expect %f = %f + %f, got %f.\n", i, expect, floatbuf[i], floatbuf[N+i], zev1[i])
		}
	}

	if fail {
		os.Exit(1)
	}

	fmt.Println("Test Passed!!!")
}
