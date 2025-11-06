import os from 'os';
import {
  complexCompute,
  ComputeResult,
} from './compute';
import {
  MemoryTracker,
  CPUMonitor,
  printCPUInfo,
  printMemoryStats,
} from './metrics';

async function parseArgs(): Promise<{ limit: number; fib: number }> {
  let limit = 100000;
  let fib = 35;

  for (let i = 2; i < process.argv.length; i++) {
    if (process.argv[i] === '--limit' && i + 1 < process.argv.length) {
      limit = parseInt(process.argv[i + 1], 10);
      i++;
    } else if (process.argv[i] === '--fib' && i + 1 < process.argv.length) {
      fib = parseInt(process.argv[i + 1], 10);
      i++;
    }
  }

  return { limit, fib };
}

async function main(): Promise<void> {
  const { limit, fib } = await parseArgs();

  // Print CPU info
  printCPUInfo();
  console.log();

  // Get number of cores
  const numCores = os.cpus().length;

  // Initialize memory tracker
  const memTracker = new MemoryTracker();
  memTracker.start();

  // Initialize CPU monitor
  const cpuMonitor = new CPUMonitor(numCores);
  cpuMonitor.startMonitoring(300);

  // Run benchmark
  const startTime = Date.now();
  const results = await complexCompute(limit, fib, numCores);
  const elapsed = (Date.now() - startTime) / 1000;

  // Stop CPU monitoring
  cpuMonitor.stopMonitoring();

  // Get final memory stats
  const startMem = memTracker.getStart();
  const peakMem = memTracker.getPeak();
  memTracker.update();
  const endMem = memTracker.getCurrent();

  // Get CPU workload distribution
  const [workDist, totalWork] = cpuMonitor.getWorkDistribution();
  const avgUtil = cpuMonitor.getAverageUtilization();

  // Print structured output
  console.log('--- TypeScript Benchmark ---');
  console.log(`Cores Used: ${numCores}`);
  console.log(`Primes Counted To: ${limit}`);
  console.log(`Fibonacci Limit: ${fib}`);
  console.log(`Primes Found: ${results.primes}`);
  console.log(`Fibonacci Sum: ${results.fibonacci}`);
  console.log(`Time Elapsed: ${elapsed.toFixed(2)}s`);
  console.log(`Memory Start: ${startMem.heapUsedMB.toFixed(2)} MB`);
  console.log(`Memory Peak: ${peakMem.heapUsedMB.toFixed(2)} MB`);
  console.log(`Memory End: ${endMem.heapUsedMB.toFixed(2)} MB`);
  console.log(`Total System Memory: ${startMem.totalGB.toFixed(2)} GB`);
  console.log(`System Memory Usage: ${startMem.usagePercent.toFixed(1)}%`);

  // Print per-core distribution
  console.log('\nPer-Core Workload Distribution:');
  const avgWorkPerCore = totalWork / numCores;
  for (let i = 0; i < numCores; i++) {
    const work = workDist[i] || 0;
    const percentage = totalWork > 0 ? (work / totalWork) * 100 : 0;
    const utilization = avgUtil[i] || 0;
    console.log(
      `  Core ${i}: ${work} units (${percentage.toFixed(1)}%) | Avg Utilization: ${utilization.toFixed(1)}%`
    );
  }
  console.log(`  Average Work per Core: ${avgWorkPerCore.toFixed(0)} units`);
}

main().catch(console.error);
