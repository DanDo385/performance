#!/usr/bin/env python3

import argparse
import math
import os
import platform
import psutil
import time
from multiprocessing import Pool, cpu_count, Manager
from dataclasses import dataclass


def solve_lattice_paths(m, n):
    """Compute the number of paths in an m x n lattice using DP"""
    if m == 0 or n == 0:
        return 0

    # Limit dimensions to avoid overflow
    m = min(m, 50)
    n = min(n, 50)

    MOD = 1000000007  # Use modulo to prevent overflow

    # DP table for lattice paths
    dp = [[0] * (n + 1) for _ in range(m + 1)]
    for i in range(m + 1):
        dp[i][0] = 1
    for j in range(n + 1):
        dp[0][j] = 1

    for i in range(1, m + 1):
        for j in range(1, n + 1):
            dp[i][j] = (dp[i - 1][j] + dp[i][j - 1]) % MOD

    return dp[m][n]


def solve_knapsack(capacity, weights, values):
    """Solve the knapsack problem (NP-hard, CPU intensive)"""
    n = len(weights)
    if n == 0 or capacity <= 0:
        return 0

    # Limit size to keep computation reasonable
    n = min(n, 25)
    weights = weights[:n]
    values = values[:n]

    # DP table
    dp = [[0] * (capacity + 1) for _ in range(n + 1)]

    for i in range(1, n + 1):
        for w in range(capacity + 1):
            if weights[i - 1] <= w:
                new_val = dp[i - 1][w - weights[i - 1]] + values[i - 1]
                dp[i][w] = max(new_val, dp[i - 1][w])
            else:
                dp[i][w] = dp[i - 1][w]

    return dp[n][capacity]


def token_generating_puzzle(puzzle_id, size):
    """Solve a computational puzzle for tokens"""
    tokens = 0

    # Puzzle 1: Lattice path computation
    lattice_size = (size % 40) + 20
    tokens += solve_lattice_paths(lattice_size, lattice_size)

    # Puzzle 2: Knapsack problems
    capacity = (size % 50) + 30
    num_items = (size % 15) + 10
    weights = [(size + i * 7 + 1) % 100 for i in range(num_items)]
    values = [(size + i * 13 + 2) % 100 for i in range(num_items)]

    tokens += solve_knapsack(capacity, weights, values)

    # Puzzle 3: Additional lattice paths with different sizes
    for i in range(size % 5):
        path_size = ((size + i) % 30) + 15
        tokens += solve_lattice_paths(path_size, path_size // 2)

    return tokens


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


@dataclass
class CoreResult:
    core_id: int
    tokens_generated: int
    lattice_paths: int
    knapsack_solutions: int


def intensive_compute_per_core(args):
    """Perform token-generating puzzle computations per core"""
    limit, num_workers, core_id = args

    core_result = CoreResult(
        core_id=core_id,
        tokens_generated=0,
        lattice_paths=0,
        knapsack_solutions=0
    )

    chunk_size = limit // num_workers
    start = core_id * chunk_size
    end = start + chunk_size
    if core_id == num_workers - 1:
        end = limit

    # Process each item, solving puzzles to generate tokens
    for i in range(start, end):
        tokens = token_generating_puzzle(i, i)
        core_result.tokens_generated += tokens

        # Track lattice paths and knapsack solutions
        if i % 100 == 0:
            core_result.lattice_paths += solve_lattice_paths((i % 30) + 15, (i % 30) + 15)

        if i % 150 == 0 and i > 0:
            capacity = (i % 50) + 30
            weights = [((i * 7 + j * 11) % 100) for j in range(15)]
            values = [((i * 13 + j * 17) % 100) for j in range(15)]
            core_result.knapsack_solutions += solve_knapsack(capacity, weights, values)

    return core_result


def intensive_benchmark(limit, num_workers):
    """Run token-generating puzzle computation with per-core breakdown"""
    tasks = [(limit, num_workers, i) for i in range(num_workers)]

    with Pool(processes=num_workers) as pool:
        core_results = pool.map(intensive_compute_per_core, tasks)

    total_tokens = sum(cr.tokens_generated for cr in core_results)

    return core_results, total_tokens


def main():
    parser = argparse.ArgumentParser(description='Python CPU Benchmark')
    parser.add_argument('--limit', type=int, default=100000,
                        help='Number of puzzle iterations')
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

    # Run intensive benchmark
    start_time = time.time()
    core_results, total_tokens = intensive_benchmark(args.limit, num_cores)
    elapsed = time.time() - start_time

    # Track peak memory
    mem_end = get_memory_info(process)
    mem_peak_rss = max(mem_peak_rss, mem_end['rss_mb'])

    # Print structured output
    print("--- Python Benchmark (Token Generation via Puzzle Solving) ---")
    print(f"Cores Used: {num_cores}")
    print(f"Puzzle Iterations: {args.limit}")
    print(f"Total Tokens Generated: {total_tokens}")
    print(f"Time Elapsed: {elapsed:.2f}s")
    print(f"Memory Start: {mem_start['rss_mb']:.2f} MB")
    print(f"Memory Peak: {mem_peak_rss:.2f} MB")
    print(f"Memory End: {mem_end['rss_mb']:.2f} MB")
    print(f"Total System Memory: {mem_start['total_gb']:.2f} GB")
    print(f"System Memory Usage: {mem_start['percent_system']:.1f}%")

    # Print per-core puzzle results
    print("\nPer-Core Puzzle Results:")
    for core_res in core_results:
        print(f"  Core {core_res.core_id}:")
        print(f"    Tokens Generated: {core_res.tokens_generated}")
        print(f"    Lattice Paths: {core_res.lattice_paths}")
        print(f"    Knapsack Solutions: {core_res.knapsack_solutions}")
    print(f"  Total: {total_tokens} tokens across {num_cores} cores")
    print(f"  Average Tokens per Core: {total_tokens // num_cores}")


if __name__ == '__main__':
    main()
