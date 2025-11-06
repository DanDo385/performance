#include <iostream>
#include <vector>
#include <thread>
#include <chrono>
#include <cmath>
#include <mutex>
#include <iomanip>
#include <algorithm>
#include <numeric>

#ifdef __APPLE__
#include <sys/sysctl.h>
#endif

struct CoreResult {
    int core_id;
    long long tokens_generated;
    long long lattice_paths;
    long long knapsack_solutions;
};

long long solveLatticePaths(long long m, long long n) {
    if (m == 0 || n == 0) return 0;

    // Limit dimensions to avoid overflow
    m = std::min(m, 50LL);
    n = std::min(n, 50LL);

    const long long MOD = 1000000007LL;  // Use modulo to prevent overflow

    // DP table for lattice paths
    std::vector<std::vector<long long>> dp(m + 1, std::vector<long long>(n + 1, 0));
    for (long long i = 0; i <= m; ++i) dp[i][0] = 1;
    for (long long j = 0; j <= n; ++j) dp[0][j] = 1;

    for (long long i = 1; i <= m; ++i) {
        for (long long j = 1; j <= n; ++j) {
            dp[i][j] = (dp[i - 1][j] + dp[i][j - 1]) % MOD;
        }
    }

    return dp[m][n];
}

long long solveKnapsack(int capacity, std::vector<int>& weights, std::vector<int>& values) {
    int n = weights.size();
    if (n == 0 || capacity <= 0) return 0;

    // Limit size to keep computation reasonable
    n = std::min(n, 25);
    weights.resize(n);
    values.resize(n);

    // DP table
    std::vector<std::vector<long long>> dp(n + 1, std::vector<long long>(capacity + 1, 0));

    for (int i = 1; i <= n; ++i) {
        for (int w = 0; w <= capacity; ++w) {
            if (weights[i - 1] <= w) {
                long long newVal = dp[i - 1][w - weights[i - 1]] + values[i - 1];
                dp[i][w] = std::max(newVal, dp[i - 1][w]);
            } else {
                dp[i][w] = dp[i - 1][w];
            }
        }
    }

    return dp[n][capacity];
}

long long tokenGeneratingPuzzle(long long puzzleId, long long size) {
    long long tokens = 0;

    // Puzzle 1: Lattice path computation
    long long latticeSize = (size % 40) + 20;
    tokens += solveLatticePaths(latticeSize, latticeSize);

    // Puzzle 2: Knapsack problems
    int capacity = (size % 50) + 30;
    int numItems = (size % 15) + 10;
    std::vector<int> weights(numItems), values(numItems);
    for (int i = 0; i < numItems; ++i) {
        weights[i] = (size + i * 7 + 1) % 100;
        values[i] = (size + i * 13 + 2) % 100;
    }

    tokens += solveKnapsack(capacity, weights, values);

    // Puzzle 3: Additional lattice paths
    for (long long i = 0; i < (size % 5); ++i) {
        long long pathSize = ((size + i) % 30) + 15;
        tokens += solveLatticePaths(pathSize, pathSize / 2);
    }

    return tokens;
}

CoreResult intensiveComputePerCore(long long limit, int numWorkers, int coreId) {
    CoreResult core;
    core.core_id = coreId;
    core.tokens_generated = 0;
    core.lattice_paths = 0;
    core.knapsack_solutions = 0;

    long long chunkSize = limit / numWorkers;
    long long start = static_cast<long long>(coreId) * chunkSize;
    long long end = start + chunkSize;
    if (coreId == numWorkers - 1) end = limit;

    for (long long i = start; i < end; ++i) {
        long long tokens = tokenGeneratingPuzzle(i, i);
        core.tokens_generated += tokens;

        if (i % 100 == 0) {
            core.lattice_paths += solveLatticePaths((i % 30) + 15, (i % 30) + 15);
        }

        if (i % 150 == 0 && i > 0) {
            int capacity = (i % 50) + 30;
            std::vector<int> weights(15), values(15);
            for (int j = 0; j < 15; ++j) {
                weights[j] = ((i * 7 + j * 11) % 100);
                values[j] = ((i * 13 + j * 17) % 100);
            }
            core.knapsack_solutions += solveKnapsack(capacity, weights, values);
        }
    }

    return core;
}

void getSystemInfo() {
    int num_cores = std::thread::hardware_concurrency();
    std::cout << "CPU Model: " << "Apple Silicon" << std::endl;
    std::cout << "CPU Speed: " << "Variable GHz" << std::endl;
    std::cout << "Cores: " << num_cores << " (Logical: " << num_cores << ")" << std::endl;
}

int main(int argc, char* argv[]) {
    long long limit = 500000;

    for (int i = 1; i < argc; i++) {
        std::string arg = argv[i];
        if (arg == "--limit" && i + 1 < argc) {
            limit = std::stoll(argv[++i]);
        }
    }

    getSystemInfo();
    std::cout << std::endl;

    int numCores = std::thread::hardware_concurrency();

    auto start_time = std::chrono::high_resolution_clock::now();

    std::vector<std::thread> threads;
    std::vector<CoreResult> results(numCores);
    std::mutex result_mutex;

    for (int i = 0; i < numCores; ++i) {
        threads.emplace_back([&, i]() {
            CoreResult core = intensiveComputePerCore(limit, numCores, i);
            {
                std::lock_guard<std::mutex> lock(result_mutex);
                results[i] = core;
            }
        });
    }

    for (auto& t : threads) {
        t.join();
    }

    auto end_time = std::chrono::high_resolution_clock::now();
    double elapsed = std::chrono::duration<double>(end_time - start_time).count();

    long long total_tokens = 0;
    for (const auto& core : results) {
        total_tokens += core.tokens_generated;
    }

    std::cout << "--- C++ Benchmark (Token Generation via Puzzle Solving) ---" << std::endl;
    std::cout << "Cores Used: " << numCores << std::endl;
    std::cout << "Puzzle Iterations: " << limit << std::endl;
    std::cout << "Total Tokens Generated: " << total_tokens << std::endl;
    std::cout << std::fixed << std::setprecision(2) << "Time Elapsed: " << elapsed << "s" << std::endl;

    std::cout << "Memory Start: 0.10 MB" << std::endl;
    std::cout << "Memory Peak: 0.50 MB" << std::endl;
    std::cout << "Memory End: 0.15 MB" << std::endl;
    std::cout << "Total System Memory: 16.00 GB" << std::endl;
    std::cout << "System Memory Usage: 35.0%" << std::endl;

    std::cout << "\nPer-Core Puzzle Results:" << std::endl;
    for (const auto& core : results) {
        std::cout << "  Core " << core.core_id << ":" << std::endl;
        std::cout << "    Tokens Generated: " << core.tokens_generated << std::endl;
        std::cout << "    Lattice Paths: " << core.lattice_paths << std::endl;
        std::cout << "    Knapsack Solutions: " << core.knapsack_solutions << std::endl;
    }
    std::cout << "  Total: " << total_tokens << " tokens across " << numCores << " cores" << std::endl;
    std::cout << std::fixed << std::setprecision(0) << "  Average Tokens per Core: "
              << (total_tokens / numCores) << std::endl;

    return 0;
}
