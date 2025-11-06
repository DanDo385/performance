#!/bin/bash

# Extract data from Go output
go_time=$(grep "Time Elapsed:" go_output.txt | awk '{print $3}' | sed 's/s//')
go_peak=$(grep "Memory Peak:" go_output.txt | awk '{print $3}')
go_cores=$(grep "Cores Used:" go_output.txt | awk '{print $3}')

# Extract data from Python output
py_time=$(grep "Time Elapsed:" python_output.txt | awk '{print $3}' | sed 's/s//')
py_peak=$(grep "Memory Peak:" python_output.txt | awk '{print $3}')
py_cores=$(grep "Cores Used:" python_output.txt | awk '{print $3}')

# Print comparison table
echo ""
echo "================================================="
echo "GO vs PYTHON CPU BENCHMARK COMPARISON"
echo "-------------------------------------------------"
printf "%-10s %-12s %-20s %-12s\n" "Language" "Time (s)" "Memory Peak (MB)" "Cores Used"
printf "%-10s %-12s %-20s %-12s\n" "Go" "$go_time" "$go_peak" "$go_cores"
printf "%-10s %-12s %-20s %-12s\n" "Python" "$py_time" "$py_peak" "$py_cores"
echo "================================================="
echo ""
