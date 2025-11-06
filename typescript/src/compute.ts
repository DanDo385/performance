import { Worker } from 'worker_threads';
import path from 'path';

export interface ComputeResult {
  primes: number;
  fibonacci: bigint;
  totalItems: number;
}

function isPrime(n: number): boolean {
  if (n <= 1) return false;
  if (n <= 3) return true;
  if (n % 2 === 0 || n % 3 === 0) return false;

  const limit = Math.sqrt(n);
  for (let i = 5; i <= limit; i += 6) {
    if (n % i === 0 || n % (i + 2) === 0) return false;
  }
  return true;
}

function countPrimesInRange(start: number, end: number): number {
  let count = 0;
  for (let i = start; i < end; i++) {
    if (isPrime(i)) count++;
  }
  return count;
}

export async function countPrimesParallel(
  limit: number,
  numWorkers: number
): Promise<number> {
  if (limit <= 0) return 0;

  const chunkSize = Math.floor(limit / numWorkers);
  const chunks: [number, number][] = [];

  for (let i = 0; i < numWorkers; i++) {
    const start = i * chunkSize;
    const end = i === numWorkers - 1 ? limit : start + chunkSize;
    chunks.push([start, end]);
  }

  // Sequential processing in TypeScript (worker threads would be for truly parallel)
  let totalCount = 0;
  for (const [start, end] of chunks) {
    totalCount += countPrimesInRange(start, end);
  }
  return totalCount;
}

function fib(n: number, memo: Map<number, bigint>): bigint {
  if (n <= 1) return BigInt(n);
  if (memo.has(n)) return memo.get(n)!;

  const result = fib(n - 1, memo) + fib(n - 2, memo);
  memo.set(n, result);
  return result;
}

export async function computeFibonacciSeries(
  numFib: number,
  numWorkers: number
): Promise<bigint> {
  if (numFib <= 0) return BigInt(0);

  const chunkSize = Math.ceil(numFib / numWorkers);
  let totalSum = BigInt(0);

  for (let i = 0; i < numWorkers; i++) {
    const start = i * chunkSize;
    const end = Math.min(start + chunkSize, numFib);

    const memo = new Map<number, bigint>();
    for (let j = start; j < end; j++) {
      totalSum += fib(j, memo);
    }
  }

  return totalSum;
}

export async function complexCompute(
  primeLimit: number,
  fibLimit: number,
  numWorkers: number
): Promise<ComputeResult> {
  const [primes, fibonacci] = await Promise.all([
    countPrimesParallel(primeLimit, numWorkers),
    computeFibonacciSeries(fibLimit, numWorkers),
  ]);

  return {
    primes,
    fibonacci,
    totalItems: primes + Number(fibonacci),
  };
}
