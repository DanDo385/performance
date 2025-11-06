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

// fib computes nth Fibonacci number with memoization for complex calculation
func fib(n int, memo map[int]int64) int64 {
	if n <= 1 {
		return int64(n)
	}
	if val, exists := memo[n]; exists {
		return val
	}

	result := fib(n-1, memo) + fib(n-2, memo)
	memo[n] = result
	return result
}

// ComputeFibonacciSeries computes Fibonacci numbers in parallel
// Returns sum of first numFib Fibonacci values
func ComputeFibonacciSeries(numFib, numWorkers int) int64 {
	if numFib <= 0 {
		return 0
	}

	// Distribute Fibonacci calculations across workers
	chunkSize := (numFib + numWorkers - 1) / numWorkers
	var wg sync.WaitGroup
	resultChan := make(chan int64, numWorkers)

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		start := i * chunkSize
		end := start + chunkSize
		if end > numFib {
			end = numFib
		}

		go func(s, e int) {
			defer wg.Done()
			memo := make(map[int]int64)
			sum := int64(0)
			for j := s; j < e; j++ {
				sum += fib(j, memo)
			}
			resultChan <- sum
		}(start, end)
	}

	// Close channel when all workers complete
	go func() {
		wg.Wait()
		close(resultChan)
	}()

	// Sum results
	totalSum := int64(0)
	for sum := range resultChan {
		totalSum += sum
	}

	return totalSum
}

// ComplexCompute performs a combined CPU-intensive computation
// mixing prime counting with Fibonacci calculations
func ComplexCompute(primeLimit, fibLimit, numWorkers int) map[string]interface{} {
	var wg sync.WaitGroup

	var primeCount int
	var fibSum int64

	wg.Add(2)

	go func() {
		defer wg.Done()
		primeCount = CountPrimesParallel(primeLimit, numWorkers)
	}()

	go func() {
		defer wg.Done()
		fibSum = ComputeFibonacciSeries(fibLimit, numWorkers)
	}()

	wg.Wait()

	return map[string]interface{}{
		"primes":     primeCount,
		"fibonacci":  fibSum,
		"totalItems": primeCount + int(fibSum),
	}
}
