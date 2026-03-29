# Quick Start — Device Enumeration

The simplest gozel example: initialize the Level Zero runtime, enumerate all available GPU drivers and their devices, and print device names.

## What It Does

- Initializes Level Zero and retrieves all GPU driver handles
- Iterates over devices under each driver, queries and prints device properties (name)

## Run

```bash
go run main.go
```

## Sample Output

```
Found 1 GPU driver(s)
  Device: Intel(R) Graphics
```
