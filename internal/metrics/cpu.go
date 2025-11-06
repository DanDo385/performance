package metrics

import (
	"fmt"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
)

// StartCPUReporter monitors and prints per-core CPU utilization
// Runs in a goroutine until stopChan is closed
func StartCPUReporter(stopChan chan struct{}) {
	ticker := time.NewTicker(300 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-stopChan:
			return
		case <-ticker.C:
			printCPUUsage()
		}
	}
}

// printCPUUsage prints current CPU utilization per core
func printCPUUsage() {
	percentages, err := cpu.Percent(0, true)
	if err != nil {
		return
	}

	fmt.Printf("\r[CPU] ")
	for i, percent := range percentages {
		fmt.Printf("Core %d: %.0f%% ", i, percent)
	}
}

// ClearCPULine clears the CPU monitoring line
func ClearCPULine() {
	fmt.Println()
}
