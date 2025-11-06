# CPU Benchmark: Multi-Language Puzzle Solver Comparison

A comprehensive cross-language CPU benchmarking project that compares performance of Go, Python, TypeScript, C, C++, and Rust on parallel computations involving complex algorithmic puzzles (lattice paths and knapsack problems).

## Features

- **Token-Generating Puzzles**: CPU-intensive algorithms combining:
  - Lattice path dynamic programming (m×n grid path enumeration)
  - NP-hard knapsack problem solving
  - Multiple puzzle iterations per benchmark run
- **6 Language Support**: Go, Python, TypeScript, C, C++, and Rust
- **Parallel Processing**: Uses all available CPU cores for computation
- **Live CPU Monitoring**: Real-time per-core CPU utilization tracking
- **Detailed CPU Info**: CPU model, speed (GHz), and core count
- **Memory Tracking**:
  - Captures start, peak, and end memory usage
  - Total system memory and usage percentage
  - Per-core computation breakdown
- **Side-by-Side Comparison**: Formatted comparison table of all languages
- **Automated Testing**: Validates benchmark output structure and values
- **Performance Analysis**: Automatic speedup calculations between languages

## Project Structure

```
performance/
├── cmd/bench/main.go                    # Go benchmark entry point
├── internal/
│   ├── compute/primes.go                # Token-generating puzzle solvers
│   ├── metrics/cpu.go                   # Live CPU monitoring & workload tracking
│   └── metrics/memory.go                # Memory tracking with system info
├── python/bench.py                      # Python benchmark implementation
├── c/src/bench.c                        # C benchmark implementation
├── cpp/src/bench.cpp                    # C++ benchmark implementation
├── rust/src/main.rs                     # Rust benchmark implementation
├── typescript/
│   ├── src/
│   │   ├── bench.ts                     # TypeScript entry point
│   │   ├── compute.ts                   # Puzzle solving logic
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

### C
- Clang or GCC compiler
- macOS: Xcode Command Line Tools
- Linux: `gcc` or `clang`

### C++
- Clang or GCC compiler with C++17 support
- macOS: Xcode Command Line Tools
- Linux: `g++` or `clang++`

### Rust
- Rust 1.56 or higher
- Cargo (comes with Rust)
- Dependencies: `rayon` for parallelism, `num_cpus` for CPU detection

### TypeScript/Node.js
- Node.js 16 or higher
- npm (comes with Node.js)
- Dependencies: TypeScript (installed via npm)

## Usage

### Run Full Benchmark (All Five Languages)

```bash
make benchmark
```

This will:
1. Build and run Go benchmark (500k puzzle iterations)
2. Run Python benchmark (same parameters)
3. Build and run TypeScript benchmark (500k puzzle iterations)
4. Build and run C benchmark (500k puzzle iterations)
5. Build and run C++ benchmark (500k puzzle iterations)
6. Build and run Rust benchmark (500k puzzle iterations)
7. Display comprehensive comparison table with performance analysis

### Run with Validation

```bash
make test
```

This runs the benchmark and validates that all outputs contain valid structured data with proper numeric values.

### Individual Language Runs

```bash
# Build and run Go
make build-go
make run-go

# Run Python
make run-python

# Build and run TypeScript
make build-ts
make run-ts

# Build and run C
make build-c
make run-c

# Build and run C++
make build-cpp
make run-cpp

# Build and run Rust
make build-rust
make run-rust
```

### Custom Parameters

```bash
# Go with custom puzzle iterations
./bench-go --limit=1000000

# Python with custom puzzle iterations
python3 python/bench.py --limit=1000000

# TypeScript with custom puzzle iterations
cd typescript && node dist/bench.js --limit=1000000

# C with custom puzzle iterations
./bench-c --limit=1000000

# C++ with custom puzzle iterations
./bench-cpp --limit=1000000

# Rust with custom puzzle iterations
./bench-rust --limit=1000000
```

## Output Format

### Individual Benchmark Output

```
CPU Model: Apple M4
CPU Speed: 2.40 GHz
Cores: 8 (Logical: 8)

--- Go Benchmark (Token Generation via Puzzle Solving) ---
Cores Used: 8
Puzzle Iterations: 500000
Total Tokens Generated: 647355583714124
Time Elapsed: 3.68s
Memory Start: 0.14 MB
Memory Peak: 1.48 MB
Memory End: 1.48 MB
Total System Memory: 16.00 GB
System Memory Usage: 75.4%

Per-Core Puzzle Results:
  Core 0:
    Tokens Generated: 80915062297108
    Lattice Paths: 233027053856
    Knapsack Solutions: 54496
    Avg Utilization: 72.4%
  Core 1:
    Tokens Generated: 80922619056391
    Lattice Paths: 233477489218
    Knapsack Solutions: 54671
    Avg Utilization: 69.0%
```

### Comparison Table

```
==========================================================================
GO vs PYTHON vs TYPESCRIPT vs C vs C++ vs RUST - CPU BENCHMARK COMPARISON
--------------------------------------------------------------------------
Language     Time (s)   Memory Peak (MB) Cores Used   Tokens Generated
C            0.95s      0.50             8            647,355,583,714,124
C++          1.01s      0.50             8            647,355,583,714,124
Rust         1.23s      0.50             8            647,355,583,714,124
Go           3.68s      1.48             8            647,355,583,714,124
Python       36.53s     15.11            8            647,355,583,714,124
==========================================================================

Performance Analysis:
--------------------
C is 3.87x faster than Go
C++ is 3.64x faster than Go
Rust is 2.99x faster than Go
Python is 9.92x slower than Go
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

### Token-Generating Puzzle
The benchmark solves three interconnected puzzles to generate computational tokens:

#### 1. Lattice Path Dynamic Programming
- **Problem**: Count number of paths from top-left to bottom-right in an m×n lattice
- **Approach**: Bottom-up DP with O(m×n) complexity
- **Computation**: For each puzzle iteration, solves multiple lattice problems with varying dimensions
- **Optimization**: Uses modulo (10^9+7) to prevent integer overflow while maintaining meaningful computation

#### 2. Knapsack Problem Solving
- **Problem**: Select items to maximize value within capacity constraint (NP-hard)
- **Approach**: Dynamic programming with O(n × capacity) complexity
- **Items**: 10-25 items per problem with random weights and values
- **Computation**: Solves multiple knapsack instances with varying parameters

#### 3. Multi-Puzzle Integration
- **Per Iteration**: Each puzzle iteration combines:
  - 1 primary lattice path solver (20-60 × 20-60 grid)
  - 1 knapsack problem
  - 0-4 additional lattice path problems with smaller grids
- **Token Value**: Sum of all puzzle results contributes to token count
- **Scaling**: 500,000 iterations × multiple puzzles per iteration = millions of computation units

## Performance Characteristics

### Actual Results (500,000 Puzzle Iterations)

| Language   | Time (s) | Tokens Generated    | Speed vs Go  | Memory Peak (MB) |
|------------|----------|---------------------|--------------|------------------|
| C          | 0.95s    | 647,355,583,714,124 | 3.87x faster | 0.50             |
| C++        | 1.01s    | 647,355,583,714,124 | 3.64x faster | 0.50             |
| Rust       | 1.23s    | 647,355,583,714,124 | 2.99x faster | 0.50             |
| Go         | 3.68s    | 647,355,583,714,124 | Baseline     | 1.48             |
| Python     | 36.53s   | 647,355,583,714,124 | 9.92x slower | 15.11            |

### Expected Relative Performance

- **C**: Fastest (compiled, minimal overhead, manual memory management)
- **C++**: Very Fast (compiled, aggressive optimization, zero-cost abstractions)
- **Rust**: Very Fast (compiled, memory safety, excellent DP optimization)
- **Go**: Fast (compiled, good parallelism, runtime overhead)
- **Python**: Slowest (interpreter overhead, GIL limitations, dynamic typing)
- **TypeScript**: Variable (JIT compiled by V8, single-threaded execution)

### Memory Usage (Typical)

- **C**: 0.5 MB (manual memory management, stack allocation)
- **C++**: 0.5 MB (static memory management, stack allocation)
- **Rust**: 0.5 MB (efficient memory management, zero-copy)
- **Go**: 1-3 MB (runtime memory management, GC overhead)
- **Python**: 10-20 MB (interpreter overhead, object allocation)
- **TypeScript**: 3-10 MB (Node.js runtime, heap management)

### CPU Utilization

- **C**: 95%+ on all cores (true parallelism with pthreads)
- **C++**: 95%+ on all cores (true parallelism with threads)
- **Rust**: 95%+ on all cores (rayon data parallelism)
- **Go**: 85%+ on all cores (goroutine overhead slightly higher)
- **Python**: 70-90% (multiprocessing IPC overhead)
- **TypeScript**: 60-80% (single-threaded, event loop scheduling)

## Testing

The validation script checks:

- All output files exist (Go, Python, TypeScript, C++, Rust)
- Required fields present: "Time Elapsed", "Memory Peak", "Cores Used", "Total Tokens Generated"
- All numeric values are valid and > 0
- Proper formatting of output
- Tokens are consistent across languages (same computation, same results)

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
- Compiled binaries (`bench-go`, `bench-c`, `bench-cpp`, `bench-rust`)
- Output files (all language output files)
- TypeScript and Rust build artifacts

## Advanced Usage

### Building Without Running

```bash
# Just compile Go
make build-go

# Just compile TypeScript
make build-ts
```

### Custom Benchmark Parameters

Edit the Makefile to change default puzzle iterations:

```makefile
run-go: build-go
    ./bench-go --limit=1000000 > go_output.txt
```

### Running Individual Benchmarks

```bash
# Run Go directly
./bench-go --limit=100000

# Run Python directly
python3 python/bench.py --limit=100000

# Run C directly
./bench-c --limit=100000

# Run C++ directly
./bench-cpp --limit=100000

# Run Rust directly
./bench-rust --limit=100000

# Run TypeScript directly
cd typescript && npm run build && node dist/bench.js --limit=100000
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

### C/C++ compiler not found

```bash
# macOS: Install Xcode Command Line Tools
xcode-select --install

# Linux: Install build essentials
sudo apt-get install build-essential
```

### Rust build fails

```bash
# Ensure Rust toolchain is installed
curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh

# Update Rust
rustup update
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
- **C**: Uses standard C library (no external deps needed)
- **C++**: Uses standard C++ library (no external deps needed)
- **Rust**: Uses standard library and Cargo dependencies
- **TypeScript**: Uses built-in `os` module (no external deps needed)

## License

MIT License - Feel free to use and modify for benchmarking purposes.

## Notes

- Results will vary based on system load and CPU governor settings
- For consistent results, run on an idle system
- Disable CPU frequency scaling for more predictable performance
- Each benchmark includes overhead from monitoring; actual computation is slightly faster
- Token counts should be identical across all languages for the same puzzle iterations and seed parameters
- The puzzle-based approach creates meaningful timing differences (no 0.00s runtimes)
- Useful for benchmarking dynamic programming and optimization algorithms across languages
