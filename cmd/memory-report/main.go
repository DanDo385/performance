package main

import (
	"fmt"
	"os"

	"github.com/benchmark/cpu-bench/internal/metrics"
)

func main() {
	if err := metrics.PrintDetailedMemoryReport(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
