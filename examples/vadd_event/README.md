# Vector Addition — Immediate Command List with Events

> [!Note]
> **SYCL** is used to write this kernel, which is not a common practice.
> Please also have a look at the **OpenCL** kernel examples like [image_scale](../image_scale/).

The same vector addition workload as the `vadd` example, but driven by an **immediate command list** and **events** instead of explicit command queues. This demonstrates fine-grained dependency tracking: memory copies signal events, and the kernel launch waits on those events before executing.

## What It Does

1. Discovers a GPU device and prints its basic & compute properties
2. Allocates host and device memory for two float32 vectors (256 MiB each)
3. Fills both vectors with random values
4. Loads a SPIR-V kernel (`vector_add`) that computes `a[i] += b[i]` in parallel
5. Creates an **event pool** with 3 events to express data-flow dependencies
6. Submits all work through a single **immediate command list**:
   - Two H→D copies, each signaling its own event
   - Kernel launch that **waits** on both copy events before executing
   - D→H copy that waits on the kernel event
7. Synchronizes via `HostSynchronize` on the immediate command list
8. Validates every element against the CPU reference

## Key Difference from `vadd`

| Aspect | `vadd` | `vadd_event` |
|--------|--------|-------------|
| Submission | 3 separate command lists executed on a command queue | 1 immediate command list |
| Synchronization | `zeCommandQueueSynchronize` | `zeCommandListHostSynchronize` |
| Dependencies | Implicit via command list ordering + barriers | Explicit via events (wait lists) |

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
Num Subgroup Sizes:          3
Subgroup Sizes:              [8 16 32 0 0 0 0 0]
=============== Computation Configuration ===============
Group Size (X, Y, Z):        (1024, 1, 1)
Group Count:                 65536
Total Elements (N):          67108864
Buffer Size:                 256 MiB
===============    Calculation Results    ===============
GPU Execution Time:          51.768500 ms
GPU Throughput:              5.19 GiB/s
===============    Validation Results    ===============
CPU Execution Time:          38.237400 ms
CPU Throughput:              7.02 GiB/s
Test Passed!!!
```
