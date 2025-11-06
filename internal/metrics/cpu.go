package metrics

import (
	"fmt"
	"sync"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
)

// CPUWorkload tracks per-core metrics and workload distribution
type CPUWorkload struct {
	mu              sync.Mutex
	coreUtilization []float64
	coreWorkUnits   []int64 // Approximate work units per core
	totalWork       int64
	samples         int
}

// NewCPUWorkload creates a new CPU workload tracker
func NewCPUWorkload(numCores int) *CPUWorkload {
	return &CPUWorkload{
		coreUtilization: make([]float64, numCores),
		coreWorkUnits:   make([]int64, numCores),
	}
}

// StartCPUReporter monitors and prints per-core CPU utilization
// Runs in a goroutine until stopChan is closed
func (cw *CPUWorkload) StartCPUReporter(stopChan chan struct{}) {
	ticker := time.NewTicker(300 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-stopChan:
			return
		case <-ticker.C:
			cw.printCPUUsage()
		}
	}
}

// printCPUUsage prints current CPU utilization per core
func (cw *CPUWorkload) printCPUUsage() {
	percentages, err := cpu.Percent(0, true)
	if err != nil {
		return
	}

	cw.mu.Lock()
	defer cw.mu.Unlock()

	// Update utilization
	for i, percent := range percentages {
		if i < len(cw.coreUtilization) {
			cw.coreUtilization[i] = percent
		}
	}
	cw.samples++

	fmt.Printf("\r[CPU] ")
	for i, percent := range percentages {
		fmt.Printf("Core %d: %.0f%% ", i, percent)
	}
}

// RecordCoreWork records work performed by a core
func (cw *CPUWorkload) RecordCoreWork(coreID int, workUnits int64) {
	cw.mu.Lock()
	defer cw.mu.Unlock()

	if coreID < len(cw.coreWorkUnits) {
		cw.coreWorkUnits[coreID] += workUnits
		cw.totalWork += workUnits
	}
}

// GetAverageUtilization returns average utilization for each core
func (cw *CPUWorkload) GetAverageUtilization() []float64 {
	cw.mu.Lock()
	defer cw.mu.Unlock()

	if cw.samples == 0 {
		return cw.coreUtilization
	}

	avg := make([]float64, len(cw.coreUtilization))
	copy(avg, cw.coreUtilization)
	return avg
}

// GetWorkDistribution returns work distribution per core
func (cw *CPUWorkload) GetWorkDistribution() ([]int64, int64) {
	cw.mu.Lock()
	defer cw.mu.Unlock()

	dist := make([]int64, len(cw.coreWorkUnits))
	copy(dist, cw.coreWorkUnits)
	return dist, cw.totalWork
}

// PrintCPUInfo prints CPU model and speed information
func PrintCPUInfo() {
	info, err := cpu.Info()
	if err != nil || len(info) == 0 {
		fmt.Println("CPU Speed: Unable to detect")
		return
	}

	cpu0 := info[0]
	fmt.Printf("CPU Model: %s\n", cpu0.ModelName)
	fmt.Printf("CPU Speed: %.2f GHz\n", cpu0.Mhz/1000)
	fmt.Printf("Cores: %d (Logical: %d)\n", cpu0.Cores, len(info))
}

// ClearCPULine clears the CPU monitoring line
func ClearCPULine() {
	fmt.Println()
}
