import os from 'os';

export interface MemoryStats {
  allocMB: number;
  totalAllocMB: number;
  heapUsedMB: number;
  heapTotalMB: number;
  externalMB: number;
  usagePercent: number;
  totalGB: number;
  freeMB: number;
}

export interface CPUInfo {
  model: string;
  speed: number; // GHz
  cores: number;
  logicalCores: number;
}

export class MemoryTracker {
  private startMem: MemoryStats;
  private peakMem: MemoryStats;
  private currentMem: MemoryStats;

  constructor() {
    this.startMem = this.getMemoryStats();
    this.peakMem = this.startMem;
    this.currentMem = this.startMem;
  }

  private getMemoryStats(): MemoryStats {
    const totalMemory = os.totalmem();
    const freeMemory = os.freemem();
    const usedMemory = totalMemory - freeMemory;

    const nodeMemUsage = process.memoryUsage();

    return {
      allocMB: nodeMemUsage.heapUsed / 1024 / 1024,
      totalAllocMB: nodeMemUsage.external / 1024 / 1024,
      heapUsedMB: nodeMemUsage.heapUsed / 1024 / 1024,
      heapTotalMB: nodeMemUsage.heapTotal / 1024 / 1024,
      externalMB: nodeMemUsage.external / 1024 / 1024,
      usagePercent: (usedMemory / totalMemory) * 100,
      totalGB: totalMemory / 1024 / 1024 / 1024,
      freeMB: freeMemory / 1024 / 1024,
    };
  }

  start(): void {
    if (global.gc) {
      global.gc();
    }
    this.startMem = this.getMemoryStats();
    this.peakMem = this.startMem;
    this.currentMem = this.startMem;
  }

  update(): void {
    this.currentMem = this.getMemoryStats();
    if (this.currentMem.allocMB > this.peakMem.allocMB) {
      this.peakMem = this.currentMem;
    }
  }

  getStart(): MemoryStats {
    return this.startMem;
  }

  getPeak(): MemoryStats {
    return this.peakMem;
  }

  getCurrent(): MemoryStats {
    this.currentMem = this.getMemoryStats();
    if (this.currentMem.allocMB > this.peakMem.allocMB) {
      this.peakMem = this.currentMem;
    }
    return this.currentMem;
  }
}

export class CPUMonitor {
  private coreUtilization: number[];
  private coreWorkUnits: number[];
  private totalWork: number = 0;
  private samples: number = 0;
  private intervalId: NodeJS.Timeout | null = null;

  constructor(numCores: number) {
    this.coreUtilization = new Array(numCores).fill(0);
    this.coreWorkUnits = new Array(numCores).fill(0);
  }

  startMonitoring(intervalMs: number = 300): void {
    const cpus = os.cpus();
    let lastCpuInfo = cpus.map((cpu) => ({
      user: cpu.times.user,
      system: cpu.times.system,
      idle: cpu.times.idle,
    }));

    this.intervalId = setInterval(() => {
      const cpus = os.cpus();
      const currentCpuInfo = cpus.map((cpu) => ({
        user: cpu.times.user,
        system: cpu.times.system,
        idle: cpu.times.idle,
      }));

      let output = '\r[CPU] ';
      for (let i = 0; i < cpus.length && i < 8; i++) {
        const lastCpu = lastCpuInfo[i];
        const currentCpu = currentCpuInfo[i];

        const userDiff = currentCpu.user - lastCpu.user;
        const systemDiff = currentCpu.system - lastCpu.system;
        const idleDiff = currentCpu.idle - lastCpu.idle;
        const total = userDiff + systemDiff + idleDiff;

        const utilization = total === 0 ? 0 : ((userDiff + systemDiff) / total) * 100;
        this.coreUtilization[i] = utilization;
        output += `Core ${i}: ${Math.round(utilization)}% `;
      }

      process.stdout.write(output);
      lastCpuInfo = currentCpuInfo;
      this.samples++;
    }, intervalMs);
  }

  stopMonitoring(): void {
    if (this.intervalId) {
      clearInterval(this.intervalId);
      this.intervalId = null;
      process.stdout.write('\n');
    }
  }

  recordCoreWork(coreId: number, workUnits: number): void {
    if (coreId < this.coreWorkUnits.length) {
      this.coreWorkUnits[coreId] += workUnits;
      this.totalWork += workUnits;
    }
  }

  getAverageUtilization(): number[] {
    return [...this.coreUtilization];
  }

  getWorkDistribution(): [number[], number] {
    return [this.coreWorkUnits, this.totalWork];
  }
}

export function getCPUInfo(): CPUInfo {
  const cpus = os.cpus();
  if (cpus.length === 0) {
    return {
      model: 'Unknown',
      speed: 0,
      cores: 0,
      logicalCores: 0,
    };
  }

  const cpu0 = cpus[0];
  return {
    model: cpu0.model,
    speed: cpu0.speed / 1000, // Convert MHz to GHz
    cores: os.cpus().length,
    logicalCores: os.cpus().length,
  };
}

export function printCPUInfo(): void {
  const info = getCPUInfo();
  console.log(`CPU Model: ${info.model}`);
  console.log(`CPU Speed: ${info.speed.toFixed(2)} GHz`);
  console.log(`Cores: ${info.cores} (Logical: ${info.logicalCores})`);
}

export function printMemoryStats(label: string, stats: MemoryStats): void {
  console.log(label);
  console.log(`  Heap Used: ${stats.heapUsedMB.toFixed(2)} MB`);
  console.log(`  Heap Total: ${stats.heapTotalMB.toFixed(2)} MB`);
  console.log(`  External: ${stats.externalMB.toFixed(2)} MB`);
  console.log(
    `  System Memory: ${((stats.totalGB - stats.freeMB / 1024).toFixed(2))} / ${stats.totalGB.toFixed(2)} GB (${stats.usagePercent.toFixed(1)}%)`
  );
}
