# CPU Benchmark: Go vs Python vs TypeScript

A comprehensive cross-language CPU benchmarking project that compares performance of Go, Python, and TypeScript on parallel computations including prime-counting and Fibonacci calculations.

## Features

- **Parallel Prime Counting**: Uses all available CPU cores for computation
- **Fibonacci Calculations**: Complex recursive calculations with memoization
- **Live CPU Monitoring**: Real-time per-core CPU utilization tracking
- **Detailed CPU Info**: CPU model, speed (GHz), and core count
- **Memory Tracking**:
  - Captures start, peak, and end memory usage
  - Total system memory and usage percentage
  - Per-core workload distribution analysis
- **Per-Core Workload Distribution**: Detailed breakdown of work units per core
- **Side-by-Side Comparison**: Formatted comparison table of all three languages
- **Automated Testing**: Validates benchmark output structure and values
- **Performance Analysis**: Automatic speedup calculations between languages

## Project Structure

```
performance/
├── cmd/bench/main.go                    # Go benchmark entry point
├── internal/
│   ├── compute/primes.go                # Parallel prime counting & Fibonacci
│   ├── metrics/cpu.go                   # Live CPU monitoring & workload tracking
│   └── metrics/memory.go                # Memory tracking with system info
├── python/bench.py                      # Python benchmark implementation
├── typescript/
│   ├── src/
│   │   ├── bench.ts                     # TypeScript entry point
│   │   ├── compute.ts                   # Computation logic
│   │   └── metrics.ts                   # Metrics collection
│   ├── package.json                     # Node.js dependencies
│   └── tsconfig.json                    # TypeScript configuration
├── go.mod                               # Go module definition
├── Makefile                             # Build and run automation
├── run_bench.sh                         # Comparison script
├── tests/validate_output.sh             # Output validation script
└── README.md                            # This file
```

## Requirements

### Go
- Go 1.21 or higher
- Dependencies (auto-installed):
  - `github.com/shirou/gopsutil/v3` - CPU and memory metrics

### Python
- Python 3.7 or higher
- Required packages:
  ```bash
  pip3 install psutil
  ```

### TypeScript/Node.js
- Node.js 16 or higher
- npm (comes with Node.js)
- Dependencies: TypeScript (installed via npm)

## Usage

### Run Full Benchmark (All Three Languages)

```bash
make benchmark
```

This will:
1. Build the Go benchmark binary
2. Run Go benchmark (500k primes, Fibonacci 35)
3. Run Python benchmark (same parameters)
4. Build and run TypeScript benchmark (same parameters)
5. Display comprehensive comparison table with performance analysis

### Run with Validation

```bash
make test
```

This runs the benchmark and validates that all outputs contain valid structured data with proper numeric values.

### Individual Language Runs

```bash
# Build Go binary
make build-go

# Run only Go benchmark
make run-go

# Run only Python benchmark
make run-python

# Build and run TypeScript
make build-ts
make run-ts
```

### Custom Parameters

```bash
# Go with custom limits
./bench-go --limit=1000000 --fib=40

# Python with custom limits
python3 python/bench.py --limit=1000000 --fib=40

# TypeScript with custom limits
cd typescript && node dist/bench.js --limit=1000000 --fib=40
```

## Output Format

### Individual Benchmark Output

```
CPU Model: Intel(R) Core(TM) i7-9700K CPU @ 3.60GHz
CPU Speed: 3.60 GHz
Cores: 8 (Logical: 8)

--- Go Benchmark ---
Cores Used: 8
Primes Counted To: 500000
Fibonacci Limit: 35
Primes Found: 41538
Fibonacci Sum: 29860703
Time Elapsed: 3.14s
Memory Start: 0.15 MB
Memory Peak: 2.45 MB
Memory End: 1.23 MB
Total System Memory: 16.00 GB
System Memory Usage: 35.4%

Per-Core Workload Distribution:
  Core 0: 5192 units (12.5%) | Avg Utilization: 95.2%
  Core 1: 5187 units (12.5%) | Avg Utilization: 94.8%
  Core 2: 5180 units (12.5%) | Avg Utilization: 93.1%
  Core 3: 5201 units (12.5%) | Avg Utilization: 96.4%
  Core 4: 5190 units (12.5%) | Avg Utilization: 94.9%
  Core 5: 5185 units (12.5%) | Avg Utilization: 95.1%
  Core 6: 5178 units (12.5%) | Avg Utilization: 92.8%
  Core 7: 5187 units (12.5%) | Avg Utilization: 94.2%
  Average Work per Core: 5188 units
```

### Comparison Table

```
==========================================================================
GO vs PYTHON vs TYPESCRIPT - CPU BENCHMARK COMPARISON
--------------------------------------------------------------------------
Language     Time (s)   Memory Peak (MB)   Cores Used   Primes     Fib Sum
Go           3.14       2.45               8            41538      29860703
Python       12.56      8.32               8            41538      29860703
TypeScript   8.42       5.67               8            41538      29860703
==========================================================================

Performance Analysis:
--------------------
Python is 4.00x slower than Go
TypeScript is 1.49x vs Python
```

## How It Works

### Go Implementation

1. **Parallel Processing**:
   - Creates worker pool using goroutines
   - Divides range across `runtime.NumCPU()` workers
   - Synchronizes using `sync.WaitGroup`

2. **CPU Monitoring**:
   - Uses `gopsutil` to read per-core utilization
   - Updates every 300ms in background goroutine
   - Displays live stats during computation

3. **Memory Tracking**:
   - Uses `runtime.ReadMemStats` for precise memory data
   - Tracks Alloc, TotalAlloc, Sys, HeapAlloc, HeapInuse
   - Captures snapshots at start, peak, and end
   - Retrieves system memory via gopsutil

4. **CPU Info**:
   - Gets CPU model name and base frequency
   - Reports number of logical cores

### Python Implementation

1. **Parallel Processing**:
   - Uses `multiprocessing.Pool` with `cpu_count()` workers
   - Splits range evenly across processes
   - Aggregates results using `pool.map()`

2. **Memory Tracking**:
   - Uses `psutil.Process().memory_info()` for RSS memory
   - Tracks virtual memory (VMS) separately
   - Gets system-wide memory stats with `psutil.virtual_memory()`

3. **CPU Info**:
   - Gets CPU model from `platform.processor()`
   - Retrieves CPU speed from `psutil.cpu_freq()` (if available)

### TypeScript Implementation

1. **Parallel Processing**:
   - Uses async/await for concurrent operations
   - Sequential computation (single-threaded JS)
   - Promise-based architecture

2. **Memory Tracking**:
   - Uses `process.memoryUsage()` for Node.js heap info
   - Gets system memory from `os.totalmem()` and `os.freemem()`
   - Calculates CPU utilization from `/proc/stat` (Linux) or native APIs

3. **CPU Info**:
   - Gets CPU model and speed from `os.cpus()`
   - Reports logical core count

## Algorithms

### Prime Counting
Both implementations use optimized trial division:
1. Handle special cases (n ≤ 1, n ≤ 3)
2. Quick checks for divisibility by 2 and 3
3. Check only numbers of form 6k±1 up to √n

### Fibonacci Sequence
- Recursive calculation with memoization
- Computes sum of first N Fibonacci numbers
- Distributed across workers for parallel execution

## Performance Characteristics

### Expected Relative Performance

- **Go**: Fastest (compiled, low-level access, true parallelism)
- **TypeScript**: Medium (V8 optimization, single-threaded JS)
- **Python**: Slowest (interpreter overhead, GIL limitations for pure Python)

### Memory Usage (Typical)

- **Go**: 0.5-3 MB (efficient memory management)
- **Python**: 5-15 MB (interpreter overhead)
- **TypeScript**: 20-50 MB (Node.js runtime)

### CPU Utilization

- **Go**: Near 100% on all cores (true parallelism)
- **Python**: Variable (GIL, multiprocessing overhead)
- **TypeScript**: Lower due to single-threaded JS execution

## Testing

The validation script checks:

- Both/all output files exist
- Required fields present: "Time Elapsed", "Memory Peak", "Cores Used", "Primes Found", "Fibonacci Sum"
- All numeric values are valid and > 0
- Proper formatting of output

```bash
make test
```

Success output:
```
✅ Benchmark output validated successfully
```

## Cleaning Up

```bash
make clean
```

Removes:
- Compiled binaries (`bench-go`)
- Output files (`go_output.txt`, `python_output.txt`, `ts_output.txt`)
- TypeScript build artifacts

## Advanced Usage

### Building Without Running

```bash
# Just compile Go
make build-go

# Just compile TypeScript
make build-ts
```

### Custom Benchmark Parameters

Edit the Makefile to change default limits:

```makefile
run-go: build-go
    ./bench-go --limit=1000000 --fib=40 > go_output.txt
```

### Running Individual Benchmarks

```bash
# Run Go directly
./bench-go --limit=100000 --fib=30

# Run Python directly
python3 python/bench.py --limit=100000 --fib=30

# Run TypeScript directly
cd typescript && npm run build && node dist/bench.js --limit=100000 --fib=30
```

## Troubleshooting

### Go build fails

```bash
# Ensure dependencies are installed
go mod tidy
go mod download
```

### Python import errors

```bash
# Install required packages
pip3 install psutil
```

### TypeScript/Node.js issues

```bash
# Clear and reinstall dependencies
rm -rf typescript/node_modules typescript/package-lock.json
cd typescript && npm install
```

### CPU/Memory detection not working

- **Go**: Requires `gopsutil` - run `go mod tidy`
- **Python**: Requires `psutil` - install with `pip3 install psutil`
- **TypeScript**: Uses built-in `os` module (no external deps needed)

## License

MIT License - Feel free to use and modify for benchmarking purposes.

## Notes

- Results will vary based on system load and CPU governor settings
- For consistent results, run on an idle system
- Disable CPU frequency scaling for more predictable performance
- Each benchmark includes overhead from monitoring; actual computation is slightly faster
