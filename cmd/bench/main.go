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
	fibLimit := flag.Int("fib", 35, "Fibonacci sequence limit")
	flag.Parse()

	// Print CPU info
	metrics.PrintCPUInfo()
	fmt.Println()

	// Configure runtime for maximum parallelism
	numCores := runtime.NumCPU()
	runtime.GOMAXPROCS(numCores)

	// Initialize memory tracker
	memTracker := metrics.NewMemoryTracker()
	memTracker.Start()

	// Initialize CPU workload tracker
	cpuWorkload := metrics.NewCPUWorkload(numCores)

	// Start CPU monitoring in background
	stopCPU := make(chan struct{})
	go cpuWorkload.StartCPUReporter(stopCPU)

	// Run complex benchmark
	startTime := time.Now()
	results := compute.ComplexCompute(*limit, *fibLimit, numCores)
	elapsed := time.Since(startTime)

	// Stop CPU monitoring
	close(stopCPU)
	time.Sleep(100 * time.Millisecond) // Allow final update
	metrics.ClearCPULine()

	// Get final memory stats
	endMem := memTracker.GetCurrent()
	peakMem := memTracker.GetPeak()
	startMem := memTracker.GetStart()

	// Get CPU workload distribution
	workDist, totalWork := cpuWorkload.GetWorkDistribution()
	avgUtil := cpuWorkload.GetAverageUtilization()

	// Print structured output
	fmt.Println("--- Go Benchmark ---")
	fmt.Printf("Cores Used: %d\n", numCores)
	fmt.Printf("Primes Counted To: %d\n", *limit)
	fmt.Printf("Fibonacci Limit: %d\n", *fibLimit)
	fmt.Printf("Primes Found: %d\n", results["primes"])
	fmt.Printf("Fibonacci Sum: %d\n", results["fibonacci"])
	fmt.Printf("Time Elapsed: %.2fs\n", elapsed.Seconds())
	fmt.Printf("Memory Start: %.2f MB\n", startMem.AllocMB)
	fmt.Printf("Memory Peak: %.2f MB\n", peakMem.AllocMB)
	fmt.Printf("Memory End: %.2f MB\n", endMem.AllocMB)
	fmt.Printf("Total System Memory: %.2f GB\n", startMem.TotalGB)
	fmt.Printf("System Memory Usage: %.1f%%\n", startMem.UsagePercent)

	// Print per-core distribution
	fmt.Println("\nPer-Core Workload Distribution:")
	avgWorkPerCore := float64(totalWork) / float64(numCores)
	for i, work := range workDist {
		percentage := 0.0
		if totalWork > 0 {
			percentage = (float64(work) / float64(totalWork)) * 100
		}
		utilization := 0.0
		if i < len(avgUtil) {
			utilization = avgUtil[i]
		}
		fmt.Printf("  Core %d: %d units (%.1f%%) | Avg Utilization: %.1f%%\n",
			i, work, percentage, utilization)
	}
	fmt.Printf("  Average Work per Core: %.0f units\n", avgWorkPerCore)
}
