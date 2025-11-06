package compute

import (
	"math"
	"sync"
)

// isPrime checks if a number is prime using trial division
func isPrime(n int) bool {
	if n <= 1 {
		return false
	}
	if n <= 3 {
		return true
	}
	if n%2 == 0 || n%3 == 0 {
		return false
	}

	limit := int(math.Sqrt(float64(n)))
	for i := 5; i <= limit; i += 6 {
		if n%i == 0 || n%(i+2) == 0 {
			return false
		}
	}
	return true
}

// countPrimesInRange counts primes in the range [start, end)
func countPrimesInRange(start, end int) int {
	count := 0
	for i := start; i < end; i++ {
		if isPrime(i) {
			count++
		}
	}
	return count
}

// CountPrimesParallel counts primes up to limit using parallel workers
func CountPrimesParallel(limit, numWorkers int) int {
	if limit <= 0 {
		return 0
	}

	chunkSize := limit / numWorkers
	var wg sync.WaitGroup
	resultChan := make(chan int, numWorkers)

	// Launch workers
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		start := i * chunkSize
		end := start + chunkSize

		// Last worker handles remainder
		if i == numWorkers-1 {
			end = limit
		}

		go func(s, e int) {
			defer wg.Done()
			count := countPrimesInRange(s, e)
			resultChan <- count
		}(start, end)
	}

	// Close channel when all workers complete
	go func() {
		wg.Wait()
		close(resultChan)
	}()

	// Sum results
	totalCount := 0
	for count := range resultChan {
		totalCount += count
	}

	return totalCount
}
