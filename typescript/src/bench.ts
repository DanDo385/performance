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

async function parseArgs(): Promise<{ limit: number }> {
  let limit = 100000;

  for (let i = 2; i < process.argv.length; i++) {
    if (process.argv[i] === '--limit' && i + 1 < process.argv.length) {
      limit = parseInt(process.argv[i + 1], 10);
      i++;
    }
  }

  return { limit };
}

async function main(): Promise<void> {
  const { limit } = await parseArgs();

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
  const results = await complexCompute(limit, numCores);
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
  console.log('--- TypeScript Benchmark (Token Generation via Puzzle Solving) ---');
  console.log(`Cores Used: ${numCores}`);
  console.log(`Puzzle Iterations: ${limit}`);
  console.log(`Total Tokens Generated: ${results.tokensGenerated}`);
  console.log(`Time Elapsed: ${elapsed.toFixed(2)}s`);
  console.log(`Memory Start: ${startMem.heapUsedMB.toFixed(2)} MB`);
  console.log(`Memory Peak: ${peakMem.heapUsedMB.toFixed(2)} MB`);
  console.log(`Memory End: ${endMem.heapUsedMB.toFixed(2)} MB`);
  console.log(`Total System Memory: ${startMem.totalGB.toFixed(2)} GB`);
  console.log(`System Memory Usage: ${startMem.usagePercent.toFixed(1)}%`);

  // Print per-core puzzle results
  console.log('\nPer-Core Puzzle Results:');
  console.log(`  Lattice Paths Computed: ${results.latticePaths}`);
  console.log(`  Knapsack Solutions: ${results.knapsackSolutions}`);
  console.log(`  Total Tokens: ${results.tokensGenerated}`);
  console.log(`  Average Tokens per Core: ${Math.floor(results.tokensGenerated / numCores)}`);
}

main().catch(console.error);
