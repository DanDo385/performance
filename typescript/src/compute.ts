import { Worker } from 'worker_threads';
import path from 'path';

export interface ComputeResult {
  tokensGenerated: number;
  latticePaths: number;
  knapsackSolutions: number;
}

function solveLatticePaths(m: number, n: number): number {
  if (m === 0 || n === 0) return 0;

  const MOD = 1000000007;
  m = Math.min(m, 50);
  n = Math.min(n, 50);

  const dp: number[][] = Array(m + 1)
    .fill(null)
    .map(() => Array(n + 1).fill(0));

  for (let i = 0; i <= m; i++) dp[i][0] = 1;
  for (let j = 0; j <= n; j++) dp[0][j] = 1;

  for (let i = 1; i <= m; i++) {
    for (let j = 1; j <= n; j++) {
      dp[i][j] = (dp[i - 1][j] + dp[i][j - 1]) % MOD;
    }
  }

  return dp[m][n];
}

function solveKnapsack(capacity: number, weights: number[], values: number[]): number {
  let n = Math.min(weights.length, 25);
  if (n === 0 || capacity <= 0) return 0;

  weights = weights.slice(0, n);
  values = values.slice(0, n);

  const dp: number[][] = Array(n + 1)
    .fill(null)
    .map(() => Array(capacity + 1).fill(0));

  for (let i = 1; i <= n; i++) {
    for (let w = 0; w <= capacity; w++) {
      if (weights[i - 1] <= w) {
        const newVal = dp[i - 1][w - weights[i - 1]] + values[i - 1];
        dp[i][w] = Math.max(newVal, dp[i - 1][w]);
      } else {
        dp[i][w] = dp[i - 1][w];
      }
    }
  }

  return dp[n][capacity];
}

function tokenGeneratingPuzzle(puzzleId: number, size: number): number {
  let tokens = 0;

  // Puzzle 1: Lattice path computation
  const latticeSize = (size % 40) + 20;
  tokens += solveLatticePaths(latticeSize, latticeSize);

  // Puzzle 2: Knapsack problems
  const capacity = (size % 50) + 30;
  const numItems = (size % 15) + 10;
  const weights: number[] = [];
  const values: number[] = [];
  for (let i = 0; i < numItems; i++) {
    weights.push((size + i * 7 + 1) % 100);
    values.push((size + i * 13 + 2) % 100);
  }

  tokens += solveKnapsack(capacity, weights, values);

  // Puzzle 3: Additional lattice paths
  for (let i = 0; i < (size % 5); i++) {
    const pathSize = ((size + i) % 30) + 15;
    tokens += solveLatticePaths(pathSize, Math.floor(pathSize / 2));
  }

  return tokens;
}

function intensiveComputePerCore(
  limit: number,
  numWorkers: number,
  coreId: number
): { coreId: number; tokensGenerated: number; latticePaths: number; knapsackSolutions: number } {
  const chunkSize = Math.floor(limit / numWorkers);
  const start = coreId * chunkSize;
  const end = coreId === numWorkers - 1 ? limit : start + chunkSize;

  let tokensGenerated = 0;
  let latticePaths = 0;
  let knapsackSolutions = 0;

  for (let i = start; i < end; i++) {
    const tokens = tokenGeneratingPuzzle(i, i);
    tokensGenerated += tokens;

    if (i % 100 === 0) {
      latticePaths += solveLatticePaths((i % 30) + 15, (i % 30) + 15);
    }

    if (i % 150 === 0 && i > 0) {
      const capacity = (i % 50) + 30;
      const weights: number[] = [];
      const values: number[] = [];
      for (let j = 0; j < 15; j++) {
        weights.push(((i * 7 + j * 11) % 100));
        values.push(((i * 13 + j * 17) % 100));
      }
      knapsackSolutions += solveKnapsack(capacity, weights, values);
    }
  }

  return { coreId, tokensGenerated, latticePaths, knapsackSolutions };
}

export async function intensiveBenchmark(
  limit: number,
  numWorkers: number
): Promise<{ results: Array<{ coreId: number; tokensGenerated: number; latticePaths: number; knapsackSolutions: number }>; totalTokens: number }> {
  const results = [];

  for (let i = 0; i < numWorkers; i++) {
    results.push(intensiveComputePerCore(limit, numWorkers, i));
  }

  const totalTokens = results.reduce((sum, r) => sum + r.tokensGenerated, 0);

  return { results, totalTokens };
}

export async function complexCompute(
  limit: number,
  numWorkers: number
): Promise<ComputeResult> {
  const { results, totalTokens } = await intensiveBenchmark(limit, numWorkers);

  const latticePaths = results.reduce((sum, r) => sum + r.latticePaths, 0);
  const knapsackSolutions = results.reduce((sum, r) => sum + r.knapsackSolutions, 0);

  return {
    tokensGenerated: totalTokens,
    latticePaths,
    knapsackSolutions,
  };
}
