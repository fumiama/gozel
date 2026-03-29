# Vector Addition — Command Queue

> ![Tip]
> **SYCL** is used to write this kernel, which is not a common practice.
> Please also have a look at the **OpenCL** kernel examples like [image_scale](../image_scale/).

A classic GPU compute example: perform element-wise addition of two large float32 vectors on the GPU, then validate the result against a CPU reference.

## What It Does

1. Discovers a GPU device and prints its basic & compute properties
2. Allocates host and device memory for two float32 vectors (256 MiB each)
3. Fills both vectors with random values and copies them to device memory
4. Loads a SPIR-V kernel (`vector_add`) that computes `a[i] += b[i]` in parallel
5. Launches the kernel via a **command queue** with explicit command lists (pre-copy → compute → post-copy)
6. Reads back the results and validates every element against the CPU reference
7. Reports GPU vs. CPU execution time and throughput

## Run

```bash
go run main.go
```

## Sample Output

```
===============  Device Basic Properties  ===============
Running on device: ID = 32103 , Name = Intel(R) Graphics @ 0.00 GHz.
=============== Device Compute Properties ===============
Max Group Size (X, Y, Z):    (1024, 1024, 1024)
Max Group Count (X, Y, Z):   (4294967295, 4294967295, 4294967295)
Max Total Group Size:        1024
Max Shared Local Memory:     65536
Subgroup Sizes:              [8 16 32]
=============== Computation Configuration ===============
Group Size (X, Y, Z):        (1024, 1, 1)
Group Count:                 65536
Total Elements (N):          67108864
Buffer Size:                 256 MiB
===============    Calculation Results    ===============
GPU Execution Time:          53.858600 ms
GPU Throughput:              4.98 GiB/s
===============    Validation Results    ===============
CPU Execution Time:          65.882900 ms
CPU Throughput:              4.07 GiB/s
Test Passed!!!
```
