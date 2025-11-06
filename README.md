# CPU Benchmark: Go vs Python

A cross-language CPU benchmarking project that compares Go and Python performance on parallel prime-counting computations.

## Features

- **Parallel Prime Counting**: Uses all available CPU cores for computation
- **Live CPU Monitoring**: Real-time per-core CPU utilization tracking (Go)
- **Memory Tracking**: Captures start, peak, and end memory usage
- **Side-by-Side Comparison**: Formatted comparison table of results
- **Automated Testing**: Validates benchmark output structure and values

## Project Structure

```
cpu-bench/
├── cmd/bench/main.go          # Go benchmark entry point
├── internal/
│   ├── compute/primes.go      # Parallel prime counting logic
│   ├── metrics/cpu.go         # Live CPU monitoring
│   └── metrics/memory.go      # Memory tracking utilities
├── python/bench.py            # Python benchmark implementation
├── tests/validate_output.sh   # Output validation tests
├── Makefile                   # Build and run automation
└── run_bench.sh              # Comparison script
```

## Requirements

### Go
- Go 1.21 or higher
- Dependencies (auto-installed):
  - `github.com/shirou/gopsutil/v3`

### Python
- Python 3.7 or higher
- Required packages:
  ```bash
  pip3 install psutil
  ```

## Usage

### Run Full Benchmark

```bash
make benchmark
```

This will:
1. Build the Go benchmark binary
2. Run Go benchmark (counting primes up to 500,000)
3. Run Python benchmark (same limit)
4. Display side-by-side comparison

### Run with Validation

```bash
make test
```

This runs the benchmark and validates that both outputs contain valid structured data.

### Individual Runs

```bash
# Build Go binary
make build-go

# Run only Go benchmark
make run-go

# Run only Python benchmark
make run-python
```

### Custom Prime Limit

```bash
# Go
./bench-go --limit=1000000

# Python
python3 python/bench.py --limit=1000000
```

## Output Format

### Individual Benchmark Output

```
--- Go Benchmark ---
Cores Used: 8
Primes Counted To: 500000
Primes Found: 41538
Time Elapsed: 3.14s
Memory Start: 0.15 MB
Memory Peak: 2.45 MB
Memory End: 1.23 MB
```

### Comparison Table

```
=================================================
GO vs PYTHON CPU BENCHMARK COMPARISON
-------------------------------------------------
Language   Time (s)   Memory Peak (MB)   Cores Used
Go         3.14       2.45              8
Python     12.56      8.32              8
=================================================
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

### Python Implementation

1. **Parallel Processing**:
   - Uses `multiprocessing.Pool` with `cpu_count()` workers
   - Splits range evenly across processes
   - Aggregates results using `pool.map()`

2. **Memory Tracking**:
   - Uses `psutil.Process().memory_info()`
   - Tracks RSS (Resident Set Size) in MB
   - Captures before and after snapshots

## Prime Counting Algorithm

Both implementations use the same optimized trial division:

```
1. Handle special cases (n ≤ 1, n ≤ 3)
2. Quick checks for divisibility by 2 and 3
3. Check only numbers of form 6k±1 up to √n
```

This reduces the number of checks significantly while maintaining correctness.

## Testing

The validation script checks:

- Both output files exist
- Required fields present: "Time Elapsed", "Memory Peak", "Cores Used"
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
- Compiled binary (`bench-go`)
- Output files (`go_output.txt`, `python_output.txt`)

## Performance Notes

- **Go** typically shows 3-5x faster execution due to:
  - Compiled native code
  - Lightweight goroutines
  - Efficient memory management

- **Python** advantages:
  - Simpler syntax
  - Rich ecosystem
  - Easier debugging

- **Fair Comparison**:
  - Both use identical algorithms
  - Same number of CPU cores
  - Same input size
  - No external optimizations

## License

MIT License - Feel free to use and modify.
