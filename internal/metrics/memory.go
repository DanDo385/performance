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
