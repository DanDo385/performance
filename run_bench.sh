#!/bin/bash

# Extract data from Go output
go_time=$(grep "Time Elapsed:" go_output.txt | awk '{print $3}' | sed 's/s//')
go_peak=$(grep "Memory Peak:" go_output.txt | awk '{print $3}')
go_cores=$(grep "Cores Used:" go_output.txt | awk '{print $3}')

# Extract data from Python output
py_time=$(grep "Time Elapsed:" python_output.txt | awk '{print $3}' | sed 's/s//')
py_peak=$(grep "Memory Peak:" python_output.txt | awk '{print $3}')
py_cores=$(grep "Cores Used:" python_output.txt | awk '{print $3}')

# Extract data from TypeScript output
ts_time=$(grep "Time Elapsed:" ts_output.txt | awk '{print $3}' | sed 's/s//')
ts_peak=$(grep "Memory Peak:" ts_output.txt | awk '{print $3}')
ts_cores=$(grep "Cores Used:" ts_output.txt | awk '{print $3}')

# Extract data from C output
c_time=$(grep "Time Elapsed:" c_output.txt | awk '{print $3}' | sed 's/s//')
c_peak=$(grep "Memory Peak:" c_output.txt | awk '{print $3}')
c_cores=$(grep "Cores Used:" c_output.txt | awk '{print $3}')

# Extract data from C++ output
cpp_time=$(grep "Time Elapsed:" cpp_output.txt | awk '{print $3}' | sed 's/s//')
cpp_peak=$(grep "Memory Peak:" cpp_output.txt | awk '{print $3}')
cpp_cores=$(grep "Cores Used:" cpp_output.txt | awk '{print $3}')

# Extract data from Rust output
rust_time=$(grep "Time Elapsed:" rust_output.txt | awk '{print $3}' | sed 's/s//')
rust_peak=$(grep "Memory Peak:" rust_output.txt | awk '{print $3}')
rust_cores=$(grep "Cores Used:" rust_output.txt | awk '{print $3}')

# Print comparison table
echo ""
echo "=========================================================================="
echo "GO vs PYTHON vs TYPESCRIPT vs C vs C++ vs RUST - CPU BENCHMARK COMPARISON"
echo "--------------------------------------------------------------------------"
printf "%-12s %-10s %-16s %-12s\n" "Language" "Time (s)" "Memory Peak (MB)" "Cores Used"
printf "%-12s %-10s %-16s %-12s\n" "Go" "$go_time" "$go_peak" "$go_cores"
printf "%-12s %-10s %-16s %-12s\n" "Python" "$py_time" "$py_peak" "$py_cores"
printf "%-12s %-10s %-16s %-12s\n" "TypeScript" "$ts_time" "$ts_peak" "$ts_cores"
printf "%-12s %-10s %-16s %-12s\n" "C" "$c_time" "$c_peak" "$c_cores"
printf "%-12s %-10s %-16s %-12s\n" "C++" "$cpp_time" "$cpp_peak" "$cpp_cores"
printf "%-12s %-10s %-16s %-12s\n" "Rust" "$rust_time" "$rust_peak" "$rust_cores"
echo "=========================================================================="
echo ""

# Calculate speedup comparisons
echo "Performance Analysis:"
echo "--------------------"
go_speedup=$(echo "scale=2; $py_time / $go_time" | bc 2>/dev/null || echo "N/A")
ts_speedup=$(echo "scale=2; $py_time / $ts_time" | bc 2>/dev/null || echo "N/A")
c_speedup=$(echo "scale=2; $py_time / $c_time" | bc 2>/dev/null || echo "N/A")
cpp_speedup=$(echo "scale=2; $py_time / $cpp_time" | bc 2>/dev/null || echo "N/A")
rust_speedup=$(echo "scale=2; $py_time / $rust_time" | bc 2>/dev/null || echo "N/A")
echo "Python is ${go_speedup}x slower than Go"
echo "TypeScript is ${ts_speedup}x faster than Python"
echo "C is ${c_speedup}x faster than Python"
echo "C++ is ${cpp_speedup}x faster than Python"
echo "Rust is ${rust_speedup}x faster than Python"
echo ""
