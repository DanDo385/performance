package main

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time"

	"github.com/benchmark/cpu-bench/internal/metrics"
)

func main() {
	// Check if bench-go exists
	benchPath := "./bench-go"
	if _, err := os.Stat(benchPath); os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "Error: bench-go executable not found. Please run 'make build-go' first.\n")
		os.Exit(1)
	}

	// Get core count
	coreCount := metrics.GetCoreCount()
	if coreCount == 0 {
		fmt.Fprintf(os.Stderr, "Error: Unable to determine CPU core count\n")
		os.Exit(1)
	}

	fmt.Println("==========================================================================")
	fmt.Println("DETAILED CPU MONITORING - PER-CORE BREAKDOWN")
	fmt.Println("==========================================================================")
	fmt.Printf("Core Count: %d logical cores\n", coreCount)
	fmt.Println()

	// Create detailed CPU monitor with 250ms update interval
	monitor := metrics.NewDetailedCPUMonitor(250 * time.Millisecond)

	// Channel to stop monitoring
	stopChan := make(chan struct{})

	// Handle interrupt signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// Start CPU monitoring
	go monitor.StartDetailedMonitoring(stopChan)

	// Start the benchmark in the background
	fmt.Println("Starting benchmark...")
	fmt.Println()
	cmd := exec.Command(benchPath, "--limit=500000")
	// Redirect benchmark output to avoid interfering with CPU monitoring display
	// The benchmark output will be shown at the end
	cmd.Stdout = nil  // Will be discarded
	cmd.Stderr = os.Stderr  // Keep errors visible

	// Start the command
	if err := cmd.Start(); err != nil {
		fmt.Fprintf(os.Stderr, "Error starting benchmark: %v\n", err)
		close(stopChan)
		os.Exit(1)
	}

	// Wait for command to finish or interrupt
	done := make(chan error, 1)
	go func() {
		done <- cmd.Wait()
	}()

	// Wait for either completion or interrupt
	select {
	case err := <-done:
		if err != nil {
			fmt.Fprintf(os.Stderr, "Benchmark error: %v\n", err)
		}
	case sig := <-sigChan:
		fmt.Printf("\n\nReceived signal: %v. Stopping benchmark...\n", sig)
		cmd.Process.Kill()
		<-done
	}

	// Stop monitoring and wait a bit for final update
	close(stopChan)
	time.Sleep(300 * time.Millisecond)

	// Get final statistics
	fmt.Println()
	fmt.Println()
	fmt.Println("==========================================================================")
	fmt.Println("CPU MONITORING SUMMARY")
	fmt.Println("==========================================================================")

	coreCountStat, perCoreAvg, totalAvg := monitor.GetDetailedStats()
	if coreCountStat > 0 {
		fmt.Printf("Core Count:           %d logical cores\n", coreCountStat)
		fmt.Printf("Total System CPU:     %.2f%% (average over run duration)\n", totalAvg)
		fmt.Println()
		fmt.Println("Per-Core Average Utilization:")
		for i, avg := range perCoreAvg {
			fmt.Printf("  Core %d:            %6.2f%%\n", i, avg)
		}
	} else {
		fmt.Println("No CPU statistics collected")
	}

	fmt.Println("==========================================================================")
}

