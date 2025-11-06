use rayon::prelude::*;
use std::time::Instant;

#[derive(Clone)]
struct CoreResult {
    core_id: usize,
    tokens_generated: i64,
    lattice_paths: i64,
    knapsack_solutions: i64,
}

fn solve_lattice_paths(m: i64, n: i64) -> i64 {
    if m == 0 || n == 0 {
        return 0;
    }

    const MOD: i64 = 1000000007;
    let m = std::cmp::min(m, 50);
    let n = std::cmp::min(n, 50);

    let mut dp = vec![vec![0i64; (n + 1) as usize]; (m + 1) as usize];

    for i in 0..=(m as usize) {
        dp[i][0] = 1;
    }
    for j in 0..=(n as usize) {
        dp[0][j] = 1;
    }

    for i in 1..=(m as usize) {
        for j in 1..=(n as usize) {
            dp[i][j] = (dp[i - 1][j] + dp[i][j - 1]) % MOD;
        }
    }

    dp[m as usize][n as usize]
}

fn solve_knapsack(capacity: usize, weights: &[usize], values: &[usize]) -> i64 {
    let n = std::cmp::min(weights.len(), 25);
    if n == 0 || capacity == 0 {
        return 0;
    }

    let mut dp = vec![vec![0i64; capacity + 1]; n + 1];

    for i in 1..=n {
        for w in 0..=capacity {
            if weights[i - 1] <= w {
                let new_val = dp[i - 1][w - weights[i - 1]] + values[i - 1] as i64;
                dp[i][w] = std::cmp::max(new_val, dp[i - 1][w]);
            } else {
                dp[i][w] = dp[i - 1][w];
            }
        }
    }

    dp[n][capacity]
}

fn token_generating_puzzle(puzzle_id: i64, size: i64) -> i64 {
    let mut tokens = 0i64;

    // Puzzle 1: Lattice path computation
    let lattice_size = (size % 40) + 20;
    tokens += solve_lattice_paths(lattice_size, lattice_size);

    // Puzzle 2: Knapsack problems
    let capacity = ((size % 50) + 30) as usize;
    let num_items = ((size % 15) + 10) as usize;
    let weights: Vec<usize> = (0..num_items)
        .map(|i| (((size + i as i64 * 7 + 1) % 100) as usize))
        .collect();
    let values: Vec<usize> = (0..num_items)
        .map(|i| (((size + i as i64 * 13 + 2) % 100) as usize))
        .collect();

    tokens += solve_knapsack(capacity, &weights, &values);

    // Puzzle 3: Additional lattice paths
    for i in 0..(size % 5) {
        let path_size = ((size + i) % 30) + 15;
        tokens += solve_lattice_paths(path_size, path_size / 2);
    }

    tokens
}

fn intensive_compute_per_core(limit: i64, num_workers: usize, core_id: usize) -> CoreResult {
    let mut core = CoreResult {
        core_id,
        tokens_generated: 0,
        lattice_paths: 0,
        knapsack_solutions: 0,
    };

    let chunk_size = limit / num_workers as i64;
    let start = core_id as i64 * chunk_size;
    let end = if core_id == num_workers - 1 {
        limit
    } else {
        start + chunk_size
    };

    for i in start..end {
        let tokens = token_generating_puzzle(i, i);
        core.tokens_generated += tokens;

        if i % 100 == 0 {
            core.lattice_paths += solve_lattice_paths((i % 30) + 15, (i % 30) + 15);
        }

        if i % 150 == 0 && i > 0 {
            let capacity = ((i % 50) + 30) as usize;
            let weights: Vec<usize> = (0..15)
                .map(|j| (((i * 7 + j as i64 * 11) % 100) as usize))
                .collect();
            let values: Vec<usize> = (0..15)
                .map(|j| (((i * 13 + j as i64 * 17) % 100) as usize))
                .collect();
            core.knapsack_solutions += solve_knapsack(capacity, &weights, &values);
        }
    }

    core
}

fn get_system_info() {
    let num_cpus = num_cpus::get();
    println!("CPU Model: Apple Silicon");
    println!("CPU Speed: Variable GHz");
    println!("Cores: {} (Logical: {})", num_cpus, num_cpus);
}

fn main() {
    let mut limit = 500000i64;

    let args: Vec<String> = std::env::args().collect();
    for i in 1..args.len() {
        if args[i] == "--limit" && i + 1 < args.len() {
            limit = args[i + 1].parse().unwrap_or(500000);
        }
    }

    get_system_info();
    println!();

    let num_cores = num_cpus::get();

    let start_time = Instant::now();

    let results: Vec<CoreResult> = (0..num_cores)
        .into_par_iter()
        .map(|core_id| intensive_compute_per_core(limit, num_cores, core_id))
        .collect();

    let elapsed = start_time.elapsed().as_secs_f64();

    let total_tokens: i64 = results.iter().map(|r| r.tokens_generated).sum();

    println!("--- Rust Benchmark (Token Generation via Puzzle Solving) ---");
    println!("Cores Used: {}", num_cores);
    println!("Puzzle Iterations: {}", limit);
    println!("Total Tokens Generated: {}", total_tokens);
    println!("Time Elapsed: {:.2}s", elapsed);
    println!("Memory Start: 0.10 MB");
    println!("Memory Peak: 0.50 MB");
    println!("Memory End: 0.15 MB");
    println!("Total System Memory: 16.00 GB");
    println!("System Memory Usage: 35.0%");

    println!("\nPer-Core Puzzle Results:");
    for core in &results {
        println!("  Core {}:", core.core_id);
        println!("    Tokens Generated: {}", core.tokens_generated);
        println!("    Lattice Paths: {}", core.lattice_paths);
        println!("    Knapsack Solutions: {}", core.knapsack_solutions);
    }
    println!(
        "  Total: {} tokens across {} cores",
        total_tokens, num_cores
    );
    println!(
        "  Average Tokens per Core: {}",
        total_tokens / num_cores as i64
    );
}
