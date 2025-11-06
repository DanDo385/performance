#!/usr/bin/env python3

import argparse
import math
import os
import time
from multiprocessing import Pool, cpu_count

import psutil


def is_prime(n):
    """Check if a number is prime using trial division"""
    if n <= 1:
        return False
    if n <= 3:
        return True
    if n % 2 == 0 or n % 3 == 0:
        return False

    limit = int(math.sqrt(n))
    i = 5
    while i <= limit:
        if n % i == 0 or n % (i + 2) == 0:
            return False
        i += 6
    return True


def count_primes_in_range(args):
    """Count primes in the range [start, end)"""
    start, end = args
    count = 0
    for i in range(start, end):
        if is_prime(i):
            count += 1
    return count


def get_memory_mb(process):
    """Get current memory usage in MB"""
    return process.memory_info().rss / (1024 * 1024)


def count_primes_parallel(limit, num_workers):
    """Count primes up to limit using parallel workers"""
    if limit <= 0:
        return 0

    chunk_size = limit // num_workers
    ranges = []

    for i in range(num_workers):
        start = i * chunk_size
        end = start + chunk_size
        # Last worker handles remainder
        if i == num_workers - 1:
            end = limit
        ranges.append((start, end))

    with Pool(processes=num_workers) as pool:
        results = pool.map(count_primes_in_range, ranges)

    return sum(results)


def main():
    parser = argparse.ArgumentParser(description='Python CPU Benchmark')
    parser.add_argument('--limit', type=int, default=100000,
                        help='Count primes up to this limit')
    args = parser.parse_args()

    # Get CPU info
    num_cores = cpu_count()

    # Initialize memory tracking
    process = psutil.Process(os.getpid())
    mem_start = get_memory_mb(process)
    mem_peak = mem_start

    # Run benchmark
    start_time = time.time()

    # Track peak memory during computation
    prime_count = count_primes_parallel(args.limit, num_cores)

    elapsed = time.time() - start_time

    # Get final memory stats
    mem_end = get_memory_mb(process)
    mem_peak = max(mem_peak, mem_end)

    # Print structured output
    print("--- Python Benchmark ---")
    print(f"Cores Used: {num_cores}")
    print(f"Primes Counted To: {args.limit}")
    print(f"Primes Found: {prime_count}")
    print(f"Time Elapsed: {elapsed:.2f}s")
    print(f"Memory Start: {mem_start:.2f} MB")
    print(f"Memory Peak: {mem_peak:.2f} MB")
    print(f"Memory End: {mem_end:.2f} MB")


if __name__ == '__main__':
    main()
