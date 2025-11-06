package compute

import (
	"sync"
)

// CoreCalculation tracks calculations performed by a specific core
type CoreCalculation struct {
	CoreID      int
	TokensGenerated int64
	LatticePaths    int64
	KnapsackSolutions int64
}

// PerCoreResults holds per-core calculation results
type PerCoreResults struct {
	Cores       []CoreCalculation
	TotalTokens int64
}

var (
	coreResultsMutex sync.Mutex
	coreResults      map[int]*CoreCalculation
)

// solveLatticePaths computes the number of paths in an mxn lattice using DP
// This creates many dynamic programming calculations
// Uses modulo to prevent overflow while still counting paths meaningfully
func solveLatticePaths(m, n int64) int64 {
	if m == 0 || n == 0 {
		return 0
	}

	// Use smaller dimensions to avoid overflow
	if m > 50 {
		m = 50
	}
	if n > 50 {
		n = 50
	}

	const MOD = int64(1000000007) // Use modulo to prevent overflow

	// DP table for lattice paths
	dp := make([][]int64, int(m)+1)
	for i := range dp {
		dp[i] = make([]int64, int(n)+1)
		dp[i][0] = 1
	}
	for j := 0; j <= int(n); j++ {
		dp[0][j] = 1
	}

	for i := 1; i <= int(m); i++ {
		for j := 1; j <= int(n); j++ {
			dp[i][j] = (dp[i-1][j] + dp[i][j-1]) % MOD
		}
	}

	return dp[int(m)][int(n)]
}

// solveKnapsack solves the knapsack problem (NP-hard, CPU intensive)
// Returns the maximum value achievable
func solveKnapsack(capacity int, weights []int, values []int) int64 {
	n := len(weights)
	if n == 0 || capacity <= 0 {
		return 0
	}

	// Limit size to keep computation reasonable per iteration
	if n > 25 {
		n = 25
	}

	// DP table
	dp := make([][]int64, n+1)
	for i := range dp {
		dp[i] = make([]int64, capacity+1)
	}

	for i := 1; i <= n; i++ {
		for w := 0; w <= capacity; w++ {
			if weights[i-1] <= w {
				newVal := dp[i-1][w-weights[i-1]] + int64(values[i-1])
				if newVal > dp[i-1][w] {
					dp[i][w] = newVal
				} else {
					dp[i][w] = dp[i-1][w]
				}
			} else {
				dp[i][w] = dp[i-1][w]
			}
		}
	}

	return dp[n][capacity]
}

// TokenGeneratingPuzzle solves a computational puzzle for tokens
// Returns number of tokens generated
func TokenGeneratingPuzzle(puzzleID, size int64) int64 {
	tokens := int64(0)

	// Puzzle 1: Lattice path computation
	latticeSize := (size % 40) + 20 // 20-60
	tokens += solveLatticePaths(latticeSize, latticeSize)

	// Puzzle 2: Knapsack problems
	capacity := int((size % 50) + 30) // 30-80
	weights := make([]int, (size%15)+10) // 10-25 items
	values := make([]int, len(weights))

	for i := 0; i < len(weights); i++ {
		weights[i] = int((size + int64(i)*7 + 1) % 100)
		values[i] = int((size + int64(i)*13 + 2) % 100)
	}

	tokens += solveKnapsack(capacity, weights, values)

	// Puzzle 3: Additional lattice paths with different sizes
	for i := int64(0); i < (size % 5); i++ {
		pathSize := ((size + i) % 30) + 15
		tokens += solveLatticePaths(pathSize, pathSize/2)
	}

	return tokens
}

// IntensiveComputePerCore performs token-generating puzzle computations per core
func IntensiveComputePerCore(limit, numWorkers, coreID int) *CoreCalculation {
	core := &CoreCalculation{
		CoreID:      coreID,
		TokensGenerated: 0,
		LatticePaths:    0,
		KnapsackSolutions: 0,
	}

	chunkSize := limit / numWorkers
	start := coreID * chunkSize
	end := start + chunkSize
	if coreID == numWorkers-1 {
		end = limit
	}

	// Process each item in the range, solving puzzles to generate tokens
	for i := start; i < end; i++ {
		tokens := TokenGeneratingPuzzle(int64(i), int64(i))
		core.TokensGenerated += tokens

		// Track lattice paths and knapsack solutions for verification
		if i%100 == 0 {
			core.LatticePaths += solveLatticePaths(int64((i%30)+15), int64((i%30)+15))
		}
		if i%150 == 0 && i > 0 {
			capacity := ((i % 50) + 30)
			weights := make([]int, 15)
			values := make([]int, 15)
			for j := 0; j < 15; j++ {
				weights[j] = int((int64(i)*7 + int64(j)*11) % 100)
				values[j] = int((int64(i)*13 + int64(j)*17) % 100)
			}
			core.KnapsackSolutions += solveKnapsack(capacity, weights, values)
		}
	}

	return core
}

// IntensiveBenchmark runs token-generating puzzle computations with per-core breakdown
func IntensiveBenchmark(limit, numWorkers int) *PerCoreResults {
	coreResults = make(map[int]*CoreCalculation)
	var wg sync.WaitGroup
	resultsChan := make(chan *CoreCalculation, numWorkers)

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func(coreID int) {
			defer wg.Done()
			result := IntensiveComputePerCore(limit, numWorkers, coreID)
			resultsChan <- result
		}(i)
	}

	go func() {
		wg.Wait()
		close(resultsChan)
	}()

	// Collect results
	results := &PerCoreResults{
		Cores: make([]CoreCalculation, 0, numWorkers),
	}

	for core := range resultsChan {
		coreResultsMutex.Lock()
		coreResults[core.CoreID] = core
		coreResultsMutex.Unlock()

		results.Cores = append(results.Cores, *core)
		results.TotalTokens += core.TokensGenerated
	}

	return results
}
