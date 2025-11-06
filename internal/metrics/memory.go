package metrics

import (
	"fmt"
	"runtime"
	"sync"

	"github.com/shirou/gopsutil/v3/mem"
)

// MemoryStats holds memory usage information
type MemoryStats struct {
	AllocMB      float64
	TotalAllocMB float64
	SysMB        float64
	HeapAllocMB  float64
	HeapInuseMB  float64
	UsagePercent float64
	TotalGB      float64
}

// MemoryTracker tracks memory usage over time
type MemoryTracker struct {
	mu         sync.Mutex
	startMem   MemoryStats
	peakMem    MemoryStats
	currentMem MemoryStats
}

// NewMemoryTracker creates a new memory tracker
func NewMemoryTracker() *MemoryTracker {
	return &MemoryTracker{}
}

// getMemoryStats reads current memory statistics
func getMemoryStats() MemoryStats {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	// Get system memory info
	v, err := mem.VirtualMemory()
	totalGB := 0.0
	usagePercent := 0.0
	if err == nil {
		totalGB = float64(v.Total) / 1024 / 1024 / 1024
		usagePercent = v.UsedPercent
	}

	return MemoryStats{
		AllocMB:      float64(m.Alloc) / 1024 / 1024,
		TotalAllocMB: float64(m.TotalAlloc) / 1024 / 1024,
		SysMB:        float64(m.Sys) / 1024 / 1024,
		HeapAllocMB:  float64(m.HeapAlloc) / 1024 / 1024,
		HeapInuseMB:  float64(m.HeapInuse) / 1024 / 1024,
		UsagePercent: usagePercent,
		TotalGB:      totalGB,
	}
}

// Start captures the starting memory state
func (mt *MemoryTracker) Start() {
	mt.mu.Lock()
	defer mt.mu.Unlock()

	runtime.GC()
	mt.startMem = getMemoryStats()
	mt.peakMem = mt.startMem
	mt.currentMem = mt.startMem
}

// Update updates current and peak memory if higher
func (mt *MemoryTracker) Update() {
	mt.mu.Lock()
	defer mt.mu.Unlock()

	mt.currentMem = getMemoryStats()
	if mt.currentMem.AllocMB > mt.peakMem.AllocMB {
		mt.peakMem = mt.currentMem
	}
}

// GetStart returns starting memory stats
func (mt *MemoryTracker) GetStart() MemoryStats {
	mt.mu.Lock()
	defer mt.mu.Unlock()
	return mt.startMem
}

// GetPeak returns peak memory stats
func (mt *MemoryTracker) GetPeak() MemoryStats {
	mt.mu.Lock()
	defer mt.mu.Unlock()
	return mt.peakMem
}

// GetCurrent returns current memory stats
func (mt *MemoryTracker) GetCurrent() MemoryStats {
	mt.mu.Lock()
	defer mt.mu.Unlock()
	mt.currentMem = getMemoryStats()
	if mt.currentMem.AllocMB > mt.peakMem.AllocMB {
		mt.peakMem = mt.currentMem
	}
	return mt.currentMem
}

// PrintMemorySnapshot prints a formatted memory snapshot
func PrintMemorySnapshot(label string) {
	stats := getMemoryStats()
	fmt.Printf("%s\n", label)
	fmt.Printf("  Alloc: %.2f MB\n", stats.AllocMB)
	fmt.Printf("  TotalAlloc: %.2f MB\n", stats.TotalAllocMB)
	fmt.Printf("  Sys: %.2f MB\n", stats.SysMB)
	fmt.Printf("  HeapAlloc: %.2f MB\n", stats.HeapAllocMB)
	fmt.Printf("  HeapInuse: %.2f MB\n", stats.HeapInuseMB)
	if stats.TotalGB > 0 {
		fmt.Printf("  System Memory: %.2f / %.2f GB (%.1f%%)\n", stats.AllocMB/1024, stats.TotalGB, stats.UsagePercent)
	}
}

// SystemMemoryStats holds detailed system memory information
type SystemMemoryStats struct {
	TotalRAM       uint64  // Total RAM in bytes
	AvailableRAM   uint64  // Available RAM in bytes
	UsedRAM        uint64  // Used RAM in bytes
	FreeRAM        uint64  // Free RAM in bytes
	Utilization    float64 // Memory utilization percentage
	SwapTotal      uint64  // Total swap in bytes
	SwapUsed       uint64  // Used swap in bytes
	SwapFree       uint64  // Free swap in bytes
	SwapUtilization float64 // Swap utilization percentage
	RAMSpeed       string  // RAM speed (if available)
}

// GetSystemMemoryStats retrieves detailed system memory information
func GetSystemMemoryStats() (*SystemMemoryStats, error) {
	v, err := mem.VirtualMemory()
	if err != nil {
		return nil, err
	}

	stats := &SystemMemoryStats{
		TotalRAM:     v.Total,
		AvailableRAM: v.Available,
		UsedRAM:      v.Used,
		FreeRAM:      v.Free,
		Utilization:  v.UsedPercent,
	}

	// Get swap information
	s, err := mem.SwapMemory()
	if err == nil {
		stats.SwapTotal = s.Total
		stats.SwapUsed = s.Used
		stats.SwapFree = s.Free
		if s.Total > 0 {
			stats.SwapUtilization = float64(s.Used) / float64(s.Total) * 100.0
		}
	}

	// Try to get RAM speed (platform-specific)
	stats.RAMSpeed = getRAMSpeed()

	return stats, nil
}

// getRAMSpeed attempts to retrieve RAM speed from system
// Returns empty string if not available
func getRAMSpeed() string {
	// This is platform-specific and may not be available on all systems
	// For macOS, we could try sysctl or dmidecode on Linux
	// For now, return empty string as graceful fallback
	return ""
}

// PrintDetailedMemoryReport prints a comprehensive memory report
func PrintDetailedMemoryReport() error {
	stats, err := GetSystemMemoryStats()
	if err != nil {
		return fmt.Errorf("failed to retrieve memory stats: %w", err)
	}

	fmt.Println("==========================================================================")
	fmt.Println("DETAILED SYSTEM MEMORY REPORT")
	fmt.Println("==========================================================================")
	fmt.Println()

	// Convert bytes to appropriate units
	totalGB := float64(stats.TotalRAM) / 1024 / 1024 / 1024
	availableGB := float64(stats.AvailableRAM) / 1024 / 1024 / 1024
	usedGB := float64(stats.UsedRAM) / 1024 / 1024 / 1024
	freeGB := float64(stats.FreeRAM) / 1024 / 1024 / 1024

	fmt.Printf("Total RAM:       %12.2f GB  (%d bytes)\n", totalGB, stats.TotalRAM)
	fmt.Printf("Available RAM:   %12.2f GB  (%d bytes)\n", availableGB, stats.AvailableRAM)
	fmt.Printf("Used RAM:        %12.2f GB  (%d bytes)\n", usedGB, stats.UsedRAM)
	fmt.Printf("Free RAM:        %12.2f GB  (%d bytes)\n", freeGB, stats.FreeRAM)
	fmt.Printf("Memory Usage:    %12.2f%%\n", stats.Utilization)
	fmt.Println()

	// Swap information
	swapTotalGB := float64(stats.SwapTotal) / 1024 / 1024 / 1024
	if stats.SwapTotal > 0 {
		swapUsedGB := float64(stats.SwapUsed) / 1024 / 1024 / 1024
		swapFreeGB := float64(stats.SwapFree) / 1024 / 1024 / 1024
		fmt.Printf("Swap Total:      %12.2f GB  (%d bytes)\n", swapTotalGB, stats.SwapTotal)
		fmt.Printf("Swap Used:       %12.2f GB  (%d bytes)\n", swapUsedGB, stats.SwapUsed)
		fmt.Printf("Swap Free:       %12.2f GB  (%d bytes)\n", swapFreeGB, stats.SwapFree)
		fmt.Printf("Swap Usage:      %12.2f%%\n", stats.SwapUtilization)
		fmt.Println()
	} else {
		fmt.Println("Swap:            Not available or not configured")
		fmt.Println()
	}

	// RAM Speed (if available)
	if stats.RAMSpeed != "" {
		fmt.Printf("RAM Speed:       %s\n", stats.RAMSpeed)
	} else {
		fmt.Println("RAM Speed:       Not available (platform limitation)")
	}

	fmt.Println("==========================================================================")

	return nil
}
