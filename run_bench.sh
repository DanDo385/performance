#!/bin/bash

# Extract data from Go output
go_time=$(grep "Time Elapsed:" go_output.txt | awk '{print $3}' | sed 's/s//')
go_peak=$(grep "Memory Peak:" go_output.txt | awk '{print $3}')
go_cores=$(grep "Cores Used:" go_output.txt | awk '{print $3}')
go_primes=$(grep "Primes Found:" go_output.txt | awk '{print $3}')
go_fib=$(grep "Fibonacci Sum:" go_output.txt | awk '{print $3}')

# Extract data from Python output
py_time=$(grep "Time Elapsed:" python_output.txt | awk '{print $3}' | sed 's/s//')
py_peak=$(grep "Memory Peak:" python_output.txt | awk '{print $3}')
py_cores=$(grep "Cores Used:" python_output.txt | awk '{print $3}')
py_primes=$(grep "Primes Found:" python_output.txt | awk '{print $3}')
py_fib=$(grep "Fibonacci Sum:" python_output.txt | awk '{print $3}')

# Extract data from TypeScript output
ts_time=$(grep "Time Elapsed:" ts_output.txt | awk '{print $3}' | sed 's/s//')
ts_peak=$(grep "Memory Peak:" ts_output.txt | awk '{print $3}')
ts_cores=$(grep "Cores Used:" ts_output.txt | awk '{print $3}')
ts_primes=$(grep "Primes Found:" ts_output.txt | awk '{print $3}')
ts_fib=$(grep "Fibonacci Sum:" ts_output.txt | awk '{print $3}')

# Print comparison table
echo ""
echo "=========================================================================="
echo "GO vs PYTHON vs TYPESCRIPT - CPU BENCHMARK COMPARISON"
echo "--------------------------------------------------------------------------"
printf "%-12s %-10s %-16s %-12s %-10s %-10s\n" "Language" "Time (s)" "Memory Peak (MB)" "Cores Used" "Primes" "Fib Sum"
printf "%-12s %-10s %-16s %-12s %-10s %-10s\n" "Go" "$go_time" "$go_peak" "$go_cores" "$go_primes" "$go_fib"
printf "%-12s %-10s %-16s %-12s %-10s %-10s\n" "Python" "$py_time" "$py_peak" "$py_cores" "$py_primes" "$py_fib"
printf "%-12s %-10s %-16s %-12s %-10s %-10s\n" "TypeScript" "$ts_time" "$ts_peak" "$ts_cores" "$ts_primes" "$ts_fib"
echo "=========================================================================="
echo ""

# Calculate speedup comparisons
echo "Performance Analysis:"
echo "--------------------"
go_speedup=$(echo "scale=2; $py_time / $go_time" | bc 2>/dev/null || echo "N/A")
ts_speedup=$(echo "scale=2; $py_time / $ts_time" | bc 2>/dev/null || echo "N/A")
echo "Python is ${go_speedup}x slower than Go"
echo "TypeScript is ${ts_speedup}x vs Python"
echo ""
