// Package main demonstrates vector addition using the gozel Level Zero bindings.
package main

import (
	"bytes"
	_ "embed"
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"math"
	"os"
	"strconv"
	"strings"
	"time"
	"unsafe"

	_ "golang.org/x/image/webp"

	"github.com/fumiama/gozel/gozel"
	"github.com/fumiama/gozel/ze"
)

//go:generate ocloc compile -file main.cl -spv_only -options "-cl-mad-enable -cl-fast-relaxed-math -cl-finite-math-only -cl-single-precision-constant" -internal_options "-O3" -output main
//go:generate llvm-spirv -to-text main_.spv -o main.spt

//go:embed main_.spv
var kernelspv []byte

//go:embed 暖笺贺春.webp
var imagebytes []byte

func main() {
	img, format, err := image.Decode(bytes.NewReader(imagebytes))
	if err != nil {
		panic(err)
	}
	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()
	ratio := float64(width) / float64(height)
	imgrgba := image.NewRGBA(bounds)
	draw.Draw(imgrgba, bounds, img, bounds.Min, draw.Src)
	dstw, dsth := width, height
	if dstw > 512 {
		dstw = 512
		dsth = int(float64(dstw) / ratio)
	}
	if dsth > 512 {
		dsth = 512
		dstw = int(float64(dsth) * ratio)
	}
	scaleRatio := float32(float64(dstw) / float64(width))

	fmt.Println("===============   Image Information   ===============")
	fmt.Printf("%-28s %s\n", "Image Format:", format)
	fmt.Printf("%-28s %.04f\n", "Image W/H ratio:", ratio)
	fmt.Printf("%-28s %d x %d\n", "Image Size:", width, height)
	fmt.Printf("%-28s %d x %d\n", "Scale to Image Size:", dstw, dsth)
	fmt.Printf("%-28s %.04f\n", "Scale ratio:", scaleRatio)
	fmt.Printf("%-28s %d bytes\n", "Image Data Size:", len(imagebytes))

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
	name, _, _ := strings.Cut(string(prop.Name[:]), "\x00")
	fmt.Println(
		"Running on device: ID =", prop.Deviceid, ", Name =", name,
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
	fmt.Printf("%-28s %v\n", "Subgroup Sizes:", cprop.Subgroupsizes[:cprop.Numsubgroupsizes])

	mod, err := ctx.ModuleCreate(dev, kernelspv)
	if err != nil {
		panic(err)
	}
	defer mod.Destroy()

	krn, err := mod.KernelCreate("scale")
	if err != nil {
		panic(err)
	}
	defer krn.Destroy()

	gX, gY, _, err := krn.SuggestGroupSize(uint32(dstw), uint32(dsth), 1)
	if err != nil {
		panic(err)
	}

	var (
		X           = uintptr(gX)
		Y           = uintptr(gY)
		groupCountX = uint32(math.Ceil(float64(dstw) / float64(X)))
		groupCountY = uint32(math.Ceil(float64(dsth) / float64(Y)))
		srcN        = uintptr(width * height * 4)                             // 4 for RGBA
		dstN        = X * uintptr(groupCountX) * Y * uintptr(groupCountY) * 4 // 4 for RGBA
		srcbufsz    = srcN * unsafe.Sizeof(uint8(0))
		dstbufsz    = dstN * unsafe.Sizeof(uint8(0))
	)
	fmt.Println("=============== Computation Configuration ===============")
	fmt.Printf("%-28s (%d, %d, %d)\n", "Group Size (X, Y, Z):", X, Y, 1)
	fmt.Printf("%-28s (%d, %d, %d)\n", "Group Count (X, Y, Z):", groupCountX, groupCountY, 1)
	fmt.Printf("%-28s (%d, %d)\n", "Total Elements (srcN, dstN):", srcN, dstN)
	fmt.Printf("%-28s %.02f KiB\n", "Source Buffer Size:", float64(srcbufsz)/1024)
	fmt.Printf("%-28s %.02f KiB\n", "Dest Buffer Size:", float64(dstbufsz)/1024)

	q, err := ctx.CommandQueueCreate(dev, gozel.ZE_COMMAND_QUEUE_MODE_DEFAULT)
	if err != nil {
		panic(err)
	}
	defer q.Destroy()

	hbuf, err := ctx.MemAllocHost(srcbufsz, 1)
	if err != nil {
		panic(err)
	}
	defer ctx.MemFree(hbuf)

	dbuf, err := ctx.MemAllocDevice(dev, srcbufsz, 1)
	if err != nil {
		panic(err)
	}
	defer ctx.MemFree(dbuf)

	himg := unsafe.Slice((*uint8)(hbuf), srcN)
	copy(himg, imgrgba.Pix)

	rgbaFmt := gozel.ZeImageFormat{
		Layout: gozel.ZE_IMAGE_FORMAT_LAYOUT_8_8_8_8,
		Type:   gozel.ZE_IMAGE_FORMAT_TYPE_UNORM, // UNORM: bilinear sampling returns float [0,1]
		X:      gozel.ZE_IMAGE_FORMAT_SWIZZLE_R,
		Y:      gozel.ZE_IMAGE_FORMAT_SWIZZLE_G,
		Z:      gozel.ZE_IMAGE_FORMAT_SWIZZLE_B,
		W:      gozel.ZE_IMAGE_FORMAT_SWIZZLE_A,
	}
	input, err := ctx.ImageCreate(dev, 0, rgbaFmt, uint64(width), uint32(height))
	if err != nil {
		panic(err)
	}
	defer input.Destroy()

	smp, err := ctx.SamplerCreate(
		dev, gozel.ZE_SAMPLER_ADDRESS_MODE_CLAMP,
		gozel.ZE_SAMPLER_FILTER_MODE_LINEAR, 1,
	)
	if err != nil {
		panic(err)
	}
	defer smp.Destroy()

	output, err := ctx.ImageCreate(
		dev, gozel.ZE_IMAGE_FLAG_KERNEL_WRITE,
		rgbaFmt, uint64(dstw), uint32(dsth),
	)
	if err != nil {
		panic(err)
	}
	defer output.Destroy()

	err = krn.SetArgumentValue(0, input)
	if err != nil {
		panic(err)
	}
	err = krn.SetArgumentValue(1, smp)
	if err != nil {
		panic(err)
	}
	err = krn.SetArgumentValue(2, output)
	if err != nil {
		panic(err)
	}
	err = krn.SetGroupSize(uint32(X), uint32(Y), 1)
	if err != nil {
		panic(err)
	}

	lstpre, err := ctx.CommandListCreate(dev)
	if err != nil {
		panic(err)
	}
	defer lstpre.Destroy()

	err = lstpre.AppendMemoryCopy(dbuf, hbuf, srcbufsz, 0)
	if err != nil {
		panic(err)
	}
	err = lstpre.AppendBarrier(0)
	if err != nil {
		panic(err)
	}

	err = lstpre.AppendImageCopyFromMemory(input, dbuf, nil, 0)
	if err != nil {
		panic(err)
	}
	err = lstpre.AppendBarrier(0)
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
		Groupcountx: groupCountX, Groupcounty: groupCountY, Groupcountz: 1,
	}, 0)
	if err != nil {
		panic(err)
	}

	err = lstcalc.AppendBarrier(0)
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

	err = lstpost.AppendImageCopyToMemory(dbuf, output, nil, 0)
	if err != nil {
		panic(err)
	}

	err = lstpost.AppendMemoryCopy(hbuf, dbuf, dstbufsz, 0)
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
	err = q.Synchronize(math.MaxUint64)
	if err != nil {
		panic(err)
	}
	elapsed := time.Since(start)

	fmt.Println("===============    Calculation Results    ===============")
	fmt.Printf("%-28s %.6f ms\n", "GPU Execution Time:", elapsed.Seconds()*1000)
	fmt.Printf("%-28s %.2f GiB/s\n", "GPU Throughput:", float64(srcbufsz)/elapsed.Seconds()/1e9)

	newimgrgba := image.NewRGBA(image.Rect(0, 0, dstw, dsth))
	copy(newimgrgba.Pix, himg)
	file, err := os.Create("small.png")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	err = png.Encode(file, newimgrgba)
	if err != nil {
		panic(err)
	}

	fmt.Println("Test Passed!!!")
}
