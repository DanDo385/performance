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

	// Run intensive parallel benchmark with per-core tracking
	startTime := time.Now()
	intensiveResults := compute.IntensiveBenchmark(*limit, numCores)
	elapsed := time.Since(startTime)

	// Stop CPU monitoring
	close(stopCPU)
	time.Sleep(100 * time.Millisecond) // Allow final update
	metrics.ClearCPULine()

	// Get final memory stats
	endMem := memTracker.GetCurrent()
	peakMem := memTracker.GetPeak()
	startMem := memTracker.GetStart()

	// Get CPU utilization
	avgUtil := cpuWorkload.GetAverageUtilization()

	// Print structured output
	fmt.Println("--- Go Benchmark (Token Generation via Puzzle Solving) ---")
	fmt.Printf("Cores Used: %d\n", numCores)
	fmt.Printf("Puzzle Iterations: %d\n", *limit)
	fmt.Printf("Total Tokens Generated: %d\n", intensiveResults.TotalTokens)
	fmt.Printf("Time Elapsed: %.2fs\n", elapsed.Seconds())
	fmt.Printf("Memory Start: %.2f MB\n", startMem.AllocMB)
	fmt.Printf("Memory Peak: %.2f MB\n", peakMem.AllocMB)
	fmt.Printf("Memory End: %.2f MB\n", endMem.AllocMB)
	fmt.Printf("Total System Memory: %.2f GB\n", startMem.TotalGB)
	fmt.Printf("System Memory Usage: %.1f%%\n", startMem.UsagePercent)

	// Print per-core breakdown
	fmt.Println("\nPer-Core Puzzle Results:")
	for _, coreResult := range intensiveResults.Cores {
		fmt.Printf("  Core %d:\n", coreResult.CoreID)
		fmt.Printf("    Tokens Generated: %d\n", coreResult.TokensGenerated)
		fmt.Printf("    Lattice Paths: %d\n", coreResult.LatticePaths)
		fmt.Printf("    Knapsack Solutions: %d\n", coreResult.KnapsackSolutions)
		if len(avgUtil) > coreResult.CoreID {
			fmt.Printf("    Avg Utilization: %.1f%%\n", avgUtil[coreResult.CoreID])
		}
	}
	fmt.Printf("  Total: %d tokens across %d cores\n", intensiveResults.TotalTokens, numCores)
	fmt.Printf("  Average Tokens per Core: %.0f\n", float64(intensiveResults.TotalTokens)/float64(numCores))
}
