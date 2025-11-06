package main

import (
	"flag"
	"fmt"
	"runtime"
	"time"

	"github.com/benchmark/cpu-bench/internal/compute"
	"github.com/benchmark/cpu-bench/internal/metrics"
)

func main() {
	// Parse command-line arguments
	limit := flag.Int("limit", 100000, "Count primes up to this limit")
	flag.Parse()

	// Configure runtime for maximum parallelism
	numCores := runtime.NumCPU()
	runtime.GOMAXPROCS(numCores)

	// Initialize memory tracker
	memTracker := metrics.NewMemoryTracker()
	memTracker.Start()

	// Start CPU monitoring in background
	stopCPU := make(chan struct{})
	go metrics.StartCPUReporter(stopCPU)

	// Run benchmark
	startTime := time.Now()
	primeCount := compute.CountPrimesParallel(*limit, numCores)
	elapsed := time.Since(startTime)

	// Stop CPU monitoring
	close(stopCPU)
	time.Sleep(100 * time.Millisecond) // Allow final update
	metrics.ClearCPULine()

	// Get final memory stats
	endMem := memTracker.GetCurrent()
	peakMem := memTracker.GetPeak()
	startMem := memTracker.GetStart()

	// Print structured output
	fmt.Println("--- Go Benchmark ---")
	fmt.Printf("Cores Used: %d\n", numCores)
	fmt.Printf("Primes Counted To: %d\n", *limit)
	fmt.Printf("Primes Found: %d\n", primeCount)
	fmt.Printf("Time Elapsed: %.2fs\n", elapsed.Seconds())
	fmt.Printf("Memory Start: %.2f MB\n", startMem.AllocMB)
	fmt.Printf("Memory Peak: %.2f MB\n", peakMem.AllocMB)
	fmt.Printf("Memory End: %.2f MB\n", endMem.AllocMB)
}
