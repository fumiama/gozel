# Image Scaling — GPU Bilinear Resize with Sampler

Downscale an image on the GPU using Level Zero's native **image** and **sampler** objects. The sampler performs hardware-accelerated bilinear interpolation, producing a high-quality resized image in a single kernel dispatch.

## What It Does

1. Decodes an embedded WebP image (1272 × 855) and converts it to RGBA
2. Computes the target dimensions (capped at 512 px on the longest side)
3. Discovers a GPU device and prints its basic & compute properties
4. Creates a SPIR-V module from an OpenCL C kernel compiled offline
5. Uses `zeKernelSuggestGroupSize` to pick an optimal 2-D workgroup size
6. Allocates host/device memory and two Level Zero **image objects** (input & output)
7. Creates a **sampler** with clamp addressing and bilinear filtering
8. Executes three command lists via a command queue:
   - **Pre**: copy host pixels → device buffer → input image
   - **Compute**: launch the `scale` kernel
   - **Post**: copy output image → device buffer → host memory
9. Writes the result to `small.png`

## Run

```bash
go run main.go
```

## Result

| Before Scaling (1272 × 855) | After Scaling (512 × 344) |
|:----------------------------:|:-------------------------:|
| ![input](暖笺贺春.webp) | ![output](small.png) |

### Console Output

```
===============   Image Information   ===============
Image Format:                webp
Image W/H ratio:             1.4877
Image Size:                  1272 x 855
Scale to Image Size:         512 x 344
Scale ratio:                 0.4025
Image Data Size:             144802 bytes
===============  Device Basic Properties  ===============
Running on device: ID = 32103 , Name = Intel(R) Graphics @ 0.00 GHz.
=============== Device Compute Properties ===============
Max Group Size (X, Y, Z):    (1024, 1024, 1024)
Max Group Count (X, Y, Z):   (4294967295, 4294967295, 4294967295)
Max Total Group Size:        1024
Max Shared Local Memory:     65536
Subgroup Sizes:              [8 16 32]
=============== Computation Configuration ===============
Group Size (X, Y, Z):        (64, 4, 1)
Group Count (X, Y, Z):       (8, 86, 1)
Total Elements (srcN, dstN): (4350240, 704512)
Source Buffer Size:          4248.28 KiB
Dest Buffer Size:            688.00 KiB
===============    Calculation Results    ===============
GPU Execution Time:          1.579000 ms
GPU Throughput:              2.76 GiB/s
Test Passed!!!
```
