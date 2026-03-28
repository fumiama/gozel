# GoZeL

[![CI](https://github.com/fumiama/gozel/actions/workflows/ci.yml/badge.svg)](https://github.com/fumiama/gozel/actions/workflows/ci.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/fumiama/gozel.svg)](https://pkg.go.dev/github.com/fumiama/gozel)
[![License: AGPL v3](https://img.shields.io/badge/License-AGPL%20v3-blue.svg)](https://www.gnu.org/licenses/agpl-3.0)

**gozel** is a pure-Go binding for the [Intel oneAPI Level Zero](https://github.com/oneapi-src/level-zero) API — enabling direct GPU/NPU compute from Go **without cgo**.

Built on [purego](https://github.com/ebitengine/purego) and Windows syscall, gozel loads `ze_loader` at runtime via FFI, avoiding all C compiler dependencies. The entire API surface is auto-generated from the official Level Zero SDK headers, keeping bindings always in sync with upstream.

---

## Table of Contents

- [Projects Using gozel](#projects-using-gozel)
- [Features](#features)
- [Platform Support](#platform-support)
- [Quick Start](#quick-start)
- [The `ze` Package — High-Level API](#the-ze-package--high-level-api)
- [Contributing](#contributing)
- [License](#license)
- [Appendix I. Architecture](#appendix-i-architecture)
- [Appendix II. Regenerating Bindings from Headers](#appendix-ii-regenerating-bindings-from-headers)
- [Appendix III. Building SPIR-V Kernels](#appendix-iii-building-spir-v-kernels)

---

## Projects Using gozel

We maintain a list of projects built with gozel. If your project uses this library, please open a PR to add it here!

| Project | Description | Link |
|---|---|---|
| | | |

## Features

- **No cgo** — Pure Go via [purego](https://github.com/ebitengine/purego) and Windows syscall. No C toolchain required at build time.
- **Auto-generated** — All 200+ Level Zero functions generated directly from the [official C headers](https://github.com/oneapi-src/level-zero/tree/master/include), covering Core, Sysman, Tools and Runtime APIs.
- **High-level wrapper** — The `ze` sub-package provides an idiomatic Go API over the raw bindings: typed handles, error returns, automatic resource lifetime.
- **Kernel in Go** — Embed SPIR-V binaries with `//go:embed`, load and launch GPU kernels entirely from Go.
- **Extensible** — Add new examples, wrap more APIs, or plug in your own SPIR-V pipelines.

## Platform Support
> gozel targets all Intel GPUs including Intel Arc / Iris Xe / Data Center GPUs that ship a Level Zero driver. Any device visible to `ze_loader` should work.

| OS | Architecture | Runtime Requirement |
|---|---|---|
| Windows | amd64 | [Intel GPU driver](https://www.intel.com/content/www/us/en/download/785597/intel-arc-iris-xe-graphics-windows.html) with `ze_loader.dll` |
| Linux | amd64 | [Intel compute-runtime](https://github.com/intel/compute-runtime) with `libze_loader.so` |
| More | under | testing... |

## Quick Start

### Install

```bash
go get github.com/fumiama/gozel
```

### Minimal Example

```go
package main

import (
	"fmt"

	"github.com/fumiama/gozel/ze"
)

func main() {
	gpus, err := ze.InitGPUDrivers()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Found %d GPU driver(s)\n", len(gpus))

	devs, _ := gpus[0].DeviceGet()
	for _, d := range devs {
		prop, _ := d.DeviceGetProperties()
		fmt.Printf("  Device: %s\n", string(prop.Name[:]))
	}

	// Found 1 GPU driver(s)
	//   Device:  Graphics
}
```

### Vector-Add

This example adds up two large float32 vectors.

```bash
cd examples/vadd
# go generate          # Optional, requires SYCL toolchain (see below)
go run main.go
```

You will see the result like

```bash
===============  Device Basic Properties  ===============
Running on device: ID = 1 , Name = Graphics @ 4.00 GHz.
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
GPU Execution Time:          56.114500 ms
GPU Throughput:              4.78 GiB/s
===============    Validation Results    ===============
CPU Execution Time:          77.190700 ms
CPU Throughput:              3.48 GiB/s
Test Passed!!!
```

### More Examples

| Example | Description | Source |
|---|---|---|
| **vadd** | Vector addition — GPU kernel launch, memory copy, validation | [examples/vadd](examples/vadd/) |
| **vadd_event** | Vector addition with event — GPU kernel launch, memory copy, validation | [examples/vadd](examples/vadd_event/) |

## The `ze` Package — High-Level API

All examples used `ze` sub-package, which wraps the raw Level Zero bindings into idiomatic Go with typed handles and method chains:

```go
// Initialize
gpus, _ := ze.InitGPUDrivers()

// Context + Device
ctx, _ := gpus[0].ContextCreate()
devs, _ := gpus[0].DeviceGet()

// Memory
hostBuf, _ := ctx.MemAllocHost(size, alignment)
devBuf, _ := ctx.MemAllocDevice(devs[0], size, alignment)
defer ctx.MemFree(hostBuf)
defer ctx.MemFree(devBuf)

// Module + Kernel
mod, _ := ctx.ModuleCreate(devs[0], spirvBytes)
krn, _ := mod.KernelCreate("my_kernel")
krn.SetArgumentValue(0, unsafe.Sizeof(uintptr(0)), unsafe.Pointer(&devBuf))
krn.SetGroupSize(256, 1, 1)

// Command submission
q, _ := ctx.CommandQueueCreate(devs[0])
cl, _ := ctx.CommandListCreate(devs[0])
cl.AppendMemoryCopy(devBuf, hostBuf, size)
cl.AppendBarrier()
cl.AppendLaunchKernel(krn, &gozel.ZeGroupCount{Groupcountx: 4, Groupcounty: 1, Groupcountz: 1})
cl.AppendBarrier()
cl.Close()
q.ExecuteCommandLists(cl)
q.Synchronize()
```

Compared to the raw `gozel.ZeCommandListCreate(...)` calls, the `ze` package:

- Eliminates boilerplate descriptor initialization (struct types, `Stype` fields)
- Returns Go `error` instead of `ZeResult` codes
- Provides method syntax on handles (`ctx.CommandListCreate(dev)` vs standalone functions)
- Manages handle lifetimes with `defer`-friendly `Destroy()` methods

The `ze` package currently covers the most common workflows. **Contributions to expand coverage are very welcome** — see [Contributing](#contributing).



> 🙌 **Have an example to share?** We'd love to grow this table — see [Contributing](#contributing).

## Contributing

Contributions of all kinds are welcome. Some particularly impactful areas:

- **Examples** — Add new `examples/` demonstrating different GPU compute patterns (matrix multiply, reduction, image processing, etc.). Every example helps new users get started faster.
- **`ze` package coverage** — The high-level wrapper doesn't cover the full Level Zero surface yet. Wrapping additional APIs (events, fences, images, sysman queries, etc.) directly benefits all downstream users.
- **Testing** — Help improve test coverage across packages.
- **Documentation** — Improve godoc comments, add usage guides, or translate documentation.
- **Project showcase** — If you use gozel in your project, open a PR to add it to the [Projects Using gozel](#projects-using-gozel) table above.

## License

- This project is generally licensed under the [GNU Affero General Public License v3.0](LICENSE).
- The files in [gozel](gozel) folder follows their original license, which is [MIT](https://github.com/oneapi-src/level-zero/blob/master/LICENSE).

---

## Appendix I. Architecture

The FFI layer (`internal/zecall`) loads the Level Zero loader library at runtime, caches procedure addresses, and dispatches calls through `purego` or Windows syscall. All pointer arguments are protected against GC collection during syscalls via `go:uintptrescapes`.

```
gozel/
├── api.go                    # Auto-generated: registers all L0 functions at init
├── core_*.go                 # Auto-generated: Core API bindings (ze*)
├── rntm_*.go                 # Auto-generated: Runtime API bindings (zer*)
├── sysm_*.go                 # Auto-generated: Sysman API bindings (zes*)
├── tols_*.go                 # Auto-generated: Tools API bindings (zet*)
├── ze/                       # High-level idiomatic Go wrapper
│   ├── init.go               #   Driver initialization
│   ├── context.go            #   Context management
│   ├── device.go             #   Device enumeration & properties
│   ├── module.go             #   SPIR-V module loading
│   ├── kernel.go             #   Kernel creation & argument binding
│   ├── mem.go                #   Device/host memory allocation
│   └── command.go            #   Command queues, lists, barriers
├── internal/zecall/          # purego FFI layer (loads ze_loader at runtime)
├── cmd/gen/                  # Code generator: parses L0 headers → Go source
├── spec/                     # Optional L0 SDK headers for dev purpose (input to cmd/gen)
└── examples/
    ├── quick_start/          # The quick start shown in this README
    └── vadd/                 # Vector addition: SYCL kernel + Go host
```

## Appendix II. Regenerating Bindings from Headers

gozel includes a code generator (`cmd/gen`) that parses the four Level Zero API headers and produces all `core_*.go`, `rntm_*.go`, `sysm_*.go`, `tols_*.go` files plus `api.go`.

### From a local SDK

Place (or symlink) the Level Zero SDK under `spec/`:

```
spec/
├── include/level_zero/
│   ├── ze_api.h      # Core
│   ├── zer_api.h     # Runtime
│   ├── zes_api.h     # Sysman
│   └── zet_api.h     # Tools
└── lib/
```

Then run:

```bash
go run ./cmd/gen -spec ./spec
```

### From a specific release

The generator can download the SDK directly from the [level-zero releases](https://github.com/oneapi-src/level-zero/releases):

```bash
go run ./cmd/gen -spec v1.28.0
```

### Via `go generate`

```bash
go generate .
```

This invokes `cmd/gen` with the local `spec/` directory as configured in `doc.go`.

## Appendix III. Building SPIR-V Kernels

GPU kernels in the [examples](examples) folder are written in SYCL C++ and compiled to SPIR-V for embedding into Go programs, which is a little bit hacky. You can also use `ocloc`, which is a common practice and you can search for the build doc elsewhere. The build pipeline uses `go generate` directives:

```
main.cpp ──clang++ -fsycl──▶ device_kern.bc
                              │
                        sycl-post-link
                              │
                        ▼ device_kern_0.bc
                              │
                   clang++ -emit-llvm -S
                              │
                        ▼ device_kern.ll
                              │
                     llvm-spirv
                              │
                        ▼ main.spv          ← embedded via //go:embed
```

### Prerequisites

- [Intel LLVM SYCL Compiler](https://github.com/intel/llvm) (provides `clang++` with SYCL support, `sycl-post-link`, etc.)

### Build

```bash
cd examples/vadd
go generate    # compiles main.cpp → main.spv
go run main.go # runs vector addition on GPU
```
