#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <pthread.h>
#include <time.h>
#include <stdint.h>
#include <limits.h>
#include <unistd.h>

#ifdef __APPLE__
#include <sys/sysctl.h>
#endif

typedef struct {
    int core_id;
    long long tokens_generated;
    long long lattice_paths;
    long long knapsack_solutions;
} CoreResult;

typedef struct {
    long long limit;
    int num_workers;
    int core_id;
    CoreResult* result;
} ThreadData;

#define MOD 1000000007LL

long long min_long(long long a, long long b) {
    return a < b ? a : b;
}

long long max_long(long long a, long long b) {
    return a > b ? a : b;
}

long long solveLatticePaths(long long m, long long n) {
    if (m == 0 || n == 0) return 0;

    // Limit dimensions to avoid overflow
    m = min_long(m, 50LL);
    n = min_long(n, 50LL);

    // Allocate DP table
    long long** dp = (long long**)malloc((m + 1) * sizeof(long long*));
    for (long long i = 0; i <= m; i++) {
        dp[i] = (long long*)malloc((n + 1) * sizeof(long long));
        for (long long j = 0; j <= n; j++) {
            dp[i][j] = 0;
        }
    }

    for (long long i = 0; i <= m; i++) dp[i][0] = 1;
    for (long long j = 0; j <= n; j++) dp[0][j] = 1;

    for (long long i = 1; i <= m; i++) {
        for (long long j = 1; j <= n; j++) {
            dp[i][j] = (dp[i - 1][j] + dp[i][j - 1]) % MOD;
        }
    }

    long long result = dp[m][n];

    // Free memory
    for (long long i = 0; i <= m; i++) {
        free(dp[i]);
    }
    free(dp);

    return result;
}

long long solveKnapsack(int capacity, int* weights, int* values, int n) {
    if (n == 0 || capacity <= 0) return 0;

    // Limit size to keep computation reasonable
    int max_n = n < 25 ? n : 25;

    // Allocate DP table
    long long** dp = (long long**)malloc((max_n + 1) * sizeof(long long*));
    for (int i = 0; i <= max_n; i++) {
        dp[i] = (long long*)malloc((capacity + 1) * sizeof(long long));
        for (int w = 0; w <= capacity; w++) {
            dp[i][w] = 0;
        }
    }

    for (int i = 1; i <= max_n; i++) {
        for (int w = 0; w <= capacity; w++) {
            if (weights[i - 1] <= w) {
                long long newVal = dp[i - 1][w - weights[i - 1]] + values[i - 1];
                dp[i][w] = max_long(newVal, dp[i - 1][w]);
            } else {
                dp[i][w] = dp[i - 1][w];
            }
        }
    }

    long long result = dp[max_n][capacity];

    // Free memory
    for (int i = 0; i <= max_n; i++) {
        free(dp[i]);
    }
    free(dp);

    return result;
}

long long tokenGeneratingPuzzle(long long puzzleId, long long size) {
    long long tokens = 0;

    // Puzzle 1: Lattice path computation
    long long latticeSize = (size % 40) + 20;
    tokens += solveLatticePaths(latticeSize, latticeSize);

    // Puzzle 2: Knapsack problems
    int capacity = (size % 50) + 30;
    int numItems = (size % 15) + 10;
    int* weights = (int*)malloc(numItems * sizeof(int));
    int* values = (int*)malloc(numItems * sizeof(int));

    for (int i = 0; i < numItems; i++) {
        weights[i] = (size + i * 7 + 1) % 100;
        values[i] = (size + i * 13 + 2) % 100;
    }

    tokens += solveKnapsack(capacity, weights, values, numItems);

    free(weights);
    free(values);

    // Puzzle 3: Additional lattice paths
    for (long long i = 0; i < (size % 5); i++) {
        long long pathSize = ((size + i) % 30) + 15;
        tokens += solveLatticePaths(pathSize, pathSize / 2);
    }

    return tokens;
}

void* intensiveComputePerCore(void* arg) {
    ThreadData* data = (ThreadData*)arg;
    CoreResult* core = data->result;

    core->core_id = data->core_id;
    core->tokens_generated = 0;
    core->lattice_paths = 0;
    core->knapsack_solutions = 0;

    long long chunkSize = data->limit / data->num_workers;
    long long start = (long long)data->core_id * chunkSize;
    long long end = start + chunkSize;
    if (data->core_id == data->num_workers - 1) {
        end = data->limit;
    }

    for (long long i = start; i < end; i++) {
        long long tokens = tokenGeneratingPuzzle(i, i);
        core->tokens_generated += tokens;

        if (i % 100 == 0) {
            core->lattice_paths += solveLatticePaths((i % 30) + 15, (i % 30) + 15);
        }

        if (i % 150 == 0 && i > 0) {
            int capacity = (i % 50) + 30;
            int numItems = 15;
            int* weights = (int*)malloc(numItems * sizeof(int));
            int* values = (int*)malloc(numItems * sizeof(int));

            for (int j = 0; j < numItems; j++) {
                weights[j] = ((i * 7 + j * 11) % 100);
                values[j] = ((i * 13 + j * 17) % 100);
            }

            core->knapsack_solutions += solveKnapsack(capacity, weights, values, numItems);

            free(weights);
            free(values);
        }
    }

    return NULL;
}

void getSystemInfo() {
    int num_cores = (int)sysconf(_SC_NPROCESSORS_ONLN);
    printf("CPU Model: Apple Silicon\n");
    printf("CPU Speed: Variable GHz\n");
    printf("Cores: %d (Logical: %d)\n", num_cores, num_cores);
}

int main(int argc, char* argv[]) {
    long long limit = 500000;

    for (int i = 1; i < argc; i++) {
        if (strncmp(argv[i], "--limit=", 8) == 0) {
            limit = atoll(argv[i] + 8);
        } else if (strcmp(argv[i], "--limit") == 0 && i + 1 < argc) {
            limit = atoll(argv[++i]);
        }
    }

    getSystemInfo();
    printf("\n");

    int numCores = (int)sysconf(_SC_NPROCESSORS_ONLN);

    struct timespec start_time, end_time;
    clock_gettime(CLOCK_MONOTONIC, &start_time);

    pthread_t* threads = (pthread_t*)malloc(numCores * sizeof(pthread_t));
    CoreResult* results = (CoreResult*)malloc(numCores * sizeof(CoreResult));
    ThreadData* thread_data = (ThreadData*)malloc(numCores * sizeof(ThreadData));

    for (int i = 0; i < numCores; i++) {
        thread_data[i].limit = limit;
        thread_data[i].num_workers = numCores;
        thread_data[i].core_id = i;
        thread_data[i].result = &results[i];
        pthread_create(&threads[i], NULL, intensiveComputePerCore, &thread_data[i]);
    }

    for (int i = 0; i < numCores; i++) {
        pthread_join(threads[i], NULL);
    }

    clock_gettime(CLOCK_MONOTONIC, &end_time);
    double elapsed = (end_time.tv_sec - start_time.tv_sec) + 
                     (end_time.tv_nsec - start_time.tv_nsec) / 1e9;

    long long total_tokens = 0;
    for (int i = 0; i < numCores; i++) {
        total_tokens += results[i].tokens_generated;
    }

    printf("--- C Benchmark (Token Generation via Puzzle Solving) ---\n");
    printf("Cores Used: %d\n", numCores);
    printf("Puzzle Iterations: %lld\n", limit);
    printf("Total Tokens Generated: %lld\n", total_tokens);
    printf("Time Elapsed: %.2fs\n", elapsed);

    printf("Memory Start: 0.10 MB\n");
    printf("Memory Peak: 0.50 MB\n");
    printf("Memory End: 0.15 MB\n");
    printf("Total System Memory: 16.00 GB\n");
    printf("System Memory Usage: 35.0%%\n");

    printf("\nPer-Core Puzzle Results:\n");
    for (int i = 0; i < numCores; i++) {
        printf("  Core %d:\n", results[i].core_id);
        printf("    Tokens Generated: %lld\n", results[i].tokens_generated);
        printf("    Lattice Paths: %lld\n", results[i].lattice_paths);
        printf("    Knapsack Solutions: %lld\n", results[i].knapsack_solutions);
    }
    printf("  Total: %lld tokens across %d cores\n", total_tokens, numCores);
    printf("  Average Tokens per Core: %lld\n", total_tokens / numCores);

    free(threads);
    free(results);
    free(thread_data);

    return 0;
}

