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

// DetailedCPUMonitor provides detailed per-core CPU monitoring
type DetailedCPUMonitor struct {
	mu              sync.Mutex
	coreSamples     [][]float64 // Store samples for each core
	totalSamples    int
	interval        time.Duration
	updateCallback  func([]float64, float64)
}

// NewDetailedCPUMonitor creates a new detailed CPU monitor
func NewDetailedCPUMonitor(interval time.Duration) *DetailedCPUMonitor {
	return &DetailedCPUMonitor{
		interval: interval,
		coreSamples: make([][]float64, 0),
	}
}

// StartDetailedMonitoring starts detailed CPU monitoring with real-time updates
func (dm *DetailedCPUMonitor) StartDetailedMonitoring(stopChan chan struct{}) {
	ticker := time.NewTicker(dm.interval)
	defer ticker.Stop()

	// Initial sample to determine core count
	percentages, err := cpu.Percent(0, true)
	if err != nil || len(percentages) == 0 {
		return
	}

	dm.mu.Lock()
	dm.coreSamples = make([][]float64, len(percentages))
	for i := range dm.coreSamples {
		dm.coreSamples[i] = make([]float64, 0)
	}
	dm.mu.Unlock()

	for {
		select {
		case <-stopChan:
			return
		case <-ticker.C:
			dm.sampleCPU()
		}
	}
}

// sampleCPU takes a CPU sample and updates display
func (dm *DetailedCPUMonitor) sampleCPU() {
	percentages, err := cpu.Percent(0, true)
	if err != nil || len(percentages) == 0 {
		return
	}

	dm.mu.Lock()
	defer dm.mu.Unlock()

	// Ensure we have enough slots
	for len(dm.coreSamples) < len(percentages) {
		dm.coreSamples = append(dm.coreSamples, make([]float64, 0))
	}

	// Store samples
	for i, percent := range percentages {
		if i < len(dm.coreSamples) {
			dm.coreSamples[i] = append(dm.coreSamples[i], percent)
		}
	}
	dm.totalSamples++

	// Calculate total system CPU
	totalCPU := 0.0
	for _, p := range percentages {
		totalCPU += p
	}
	totalCPU /= float64(len(percentages))

	// Print detailed breakdown
	dm.printDetailedCPU(percentages, totalCPU)
}

// printDetailedCPU prints formatted CPU information
func (dm *DetailedCPUMonitor) printDetailedCPU(percentages []float64, totalCPU float64) {
	// Clear previous line and print header
	fmt.Printf("\r\033[K")
	fmt.Printf("CPU Load: Total %.1f%% | ", totalCPU)
	
	for i, percent := range percentages {
		fmt.Printf("Core %d: %5.1f%%", i, percent)
		if i < len(percentages)-1 {
			fmt.Printf(" | ")
		}
	}
}

// GetDetailedStats returns detailed CPU statistics
func (dm *DetailedCPUMonitor) GetDetailedStats() (coreCount int, perCoreAvg []float64, totalAvg float64) {
	dm.mu.Lock()
	defer dm.mu.Unlock()

	if len(dm.coreSamples) == 0 || dm.totalSamples == 0 {
		return 0, nil, 0.0
	}

	coreCount = len(dm.coreSamples)
	perCoreAvg = make([]float64, coreCount)
	
	totalSum := 0.0
	for i, samples := range dm.coreSamples {
		if len(samples) == 0 {
			perCoreAvg[i] = 0.0
			continue
		}
		
		sum := 0.0
		for _, s := range samples {
			sum += s
		}
		avg := sum / float64(len(samples))
		perCoreAvg[i] = avg
		totalSum += avg
	}
	
	totalAvg = totalSum / float64(coreCount)
	return coreCount, perCoreAvg, totalAvg
}

// GetCoreCount returns the number of CPU cores
func GetCoreCount() int {
	info, err := cpu.Info()
	if err != nil || len(info) == 0 {
		return 0
	}
	return len(info)
}
