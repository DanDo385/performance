#!/bin/bash

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
NC='\033[0m' # No Color

ERRORS=0

# Function to check if file exists
check_file_exists() {
    if [ ! -f "$1" ]; then
        echo -e "${RED}❌ Error: File $1 does not exist${NC}"
        ERRORS=$((ERRORS + 1))
        return 1
    fi
    return 0
}

# Function to check if a field exists and has a valid value
check_field() {
    local file=$1
    local field=$2
    local pattern=$3

    if ! grep -q "$field" "$file"; then
        echo -e "${RED}❌ Error: Field '$field' not found in $file${NC}"
        ERRORS=$((ERRORS + 1))
        return 1
    fi

    local value=$(grep "$field" "$file" | awk '{print $3}' | sed 's/s$//' | sed 's/MB$//')

    if ! [[ "$value" =~ $pattern ]]; then
        echo -e "${RED}❌ Error: Field '$field' in $file has invalid value: $value${NC}"
        ERRORS=$((ERRORS + 1))
        return 1
    fi

    # Check if numeric value is > 0
    if (( $(echo "$value <= 0" | bc -l 2>/dev/null || echo "0") )); then
        echo -e "${RED}❌ Error: Field '$field' in $file must be > 0, got: $value${NC}"
        ERRORS=$((ERRORS + 1))
        return 1
    fi

    return 0
}

echo "Validating benchmark outputs..."
echo ""

# Check if output files exist
check_file_exists "go_output.txt" || exit 1
check_file_exists "python_output.txt" || exit 1
check_file_exists "ts_output.txt" || exit 1

# Validate Go output
echo "Checking Go output..."
check_field "go_output.txt" "Time Elapsed:" '^[0-9]+\.?[0-9]*$'
check_field "go_output.txt" "Memory Peak:" '^[0-9]+\.?[0-9]*$'
check_field "go_output.txt" "Cores Used:" '^[0-9]+$'
check_field "go_output.txt" "Primes Found:" '^[0-9]+$'
check_field "go_output.txt" "Fibonacci Sum:" '^[0-9]+$'

# Validate Python output
echo "Checking Python output..."
check_field "python_output.txt" "Time Elapsed:" '^[0-9]+\.?[0-9]*$'
check_field "python_output.txt" "Memory Peak:" '^[0-9]+\.?[0-9]*$'
check_field "python_output.txt" "Cores Used:" '^[0-9]+$'
check_field "python_output.txt" "Primes Found:" '^[0-9]+$'
check_field "python_output.txt" "Fibonacci Sum:" '^[0-9]+$'

# Validate TypeScript output
echo "Checking TypeScript output..."
check_field "ts_output.txt" "Time Elapsed:" '^[0-9]+\.?[0-9]*$'
check_field "ts_output.txt" "Memory Peak:" '^[0-9]+\.?[0-9]*$'
check_field "ts_output.txt" "Cores Used:" '^[0-9]+$'
check_field "ts_output.txt" "Primes Found:" '^[0-9]+$'
check_field "ts_output.txt" "Fibonacci Sum:" '^[0-9]+$'

echo ""
if [ $ERRORS -eq 0 ]; then
    echo -e "${GREEN}✅ Benchmark output validated successfully${NC}"
    exit 0
else
    echo -e "${RED}❌ Validation failed with $ERRORS error(s)${NC}"
    exit 1
fi
