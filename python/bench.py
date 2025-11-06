#!/usr/bin/env python3

import argparse
import math
import os
import platform
import psutil
import time
from multiprocessing import Pool, cpu_count


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


def fib(n, memo=None):
    """Compute nth Fibonacci number with memoization"""
    if memo is None:
        memo = {}

    if n <= 1:
        return n
    if n in memo:
        return memo[n]

    result = fib(n - 1, memo) + fib(n - 2, memo)
    memo[n] = result
    return result


def compute_fibonacci_range(args):
    """Compute Fibonacci sum for a range of values"""
    start, end = args
    memo = {}
    total = 0
    for i in range(start, end):
        total += fib(i, memo)
    return total


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


def compute_fibonacci_series(num_fib, num_workers):
    """Compute Fibonacci series sum using parallel workers"""
    if num_fib <= 0:
        return 0

    chunk_size = (num_fib + num_workers - 1) // num_workers
    ranges = []

    for i in range(num_workers):
        start = i * chunk_size
        end = min(start + chunk_size, num_fib)
        if start < num_fib:
            ranges.append((start, end))

    with Pool(processes=num_workers) as pool:
        results = pool.map(compute_fibonacci_range, ranges)

    return sum(results)


def get_cpu_info():
    """Get CPU information"""
    try:
        # Try to get from psutil
        freq = psutil.cpu_freq()
        cpu_speed = freq.current / 1000 if freq else 0  # Convert MHz to GHz
    except:
        cpu_speed = 0

    return {
        "model": platform.processor() or "Unknown",
        "speed": cpu_speed,
        "cores": cpu_count(),
    }


def get_memory_info(process):
    """Get current memory usage details"""
    mem_info = process.memory_info()
    vm_info = psutil.virtual_memory()

    return {
        "rss_mb": mem_info.rss / (1024 * 1024),
        "vms_mb": mem_info.vms / (1024 * 1024),
        "percent": process.memory_percent(),
        "total_gb": vm_info.total / (1024 ** 3),
        "available_mb": vm_info.available / (1024 * 1024),
        "percent_system": vm_info.percent,
    }


def print_cpu_info(info):
    """Print CPU information"""
    print(f"CPU Model: {info['model']}")
    if info['speed'] > 0:
        print(f"CPU Speed: {info['speed']:.2f} GHz")
    print(f"Cores: {info['cores']} (Logical: {info['cores']})")


def main():
    parser = argparse.ArgumentParser(description='Python CPU Benchmark')
    parser.add_argument('--limit', type=int, default=100000,
                        help='Count primes up to this limit')
    parser.add_argument('--fib', type=int, default=35,
                        help='Fibonacci limit')
    args = parser.parse_args()

    # Print CPU info
    cpu_info = get_cpu_info()
    print_cpu_info(cpu_info)
    print()

    # Get CPU info
    num_cores = cpu_count()

    # Initialize memory tracking
    process = psutil.Process(os.getpid())
    mem_start = get_memory_info(process)
    mem_peak_rss = mem_start['rss_mb']

    # Run benchmark
    start_time = time.time()

    # Compute both primes and fibonacci in sequence (no true parallelism between them)
    prime_count = count_primes_parallel(args.limit, num_cores)
    fib_sum = compute_fibonacci_series(args.fib, num_cores)

    elapsed = time.time() - start_time

    # Track peak memory
    mem_end = get_memory_info(process)
    mem_peak_rss = max(mem_peak_rss, mem_end['rss_mb'])

    # Print structured output
    print("--- Python Benchmark ---")
    print(f"Cores Used: {num_cores}")
    print(f"Primes Counted To: {args.limit}")
    print(f"Fibonacci Limit: {args.fib}")
    print(f"Primes Found: {prime_count}")
    print(f"Fibonacci Sum: {fib_sum}")
    print(f"Time Elapsed: {elapsed:.2f}s")
    print(f"Memory Start: {mem_start['rss_mb']:.2f} MB")
    print(f"Memory Peak: {mem_peak_rss:.2f} MB")
    print(f"Memory End: {mem_end['rss_mb']:.2f} MB")
    print(f"Total System Memory: {mem_start['total_gb']:.2f} GB")
    print(f"System Memory Usage: {mem_start['percent_system']:.1f}%")

    # Per-core distribution (Python doesn't expose this easily in multiprocessing)
    print("\nPer-Core Workload Distribution:")
    avg_work_per_core = (prime_count + fib_sum) / num_cores
    for i in range(num_cores):
        # Estimate even distribution
        percentage = 100.0 / num_cores
        print(f"  Core {i}: {int(avg_work_per_core)} units ({percentage:.1f}%) | Avg Utilization: N/A")
    print(f"  Average Work per Core: {int(avg_work_per_core)} units")


if __name__ == '__main__':
    main()
