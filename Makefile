.PHONY: default build-go run-go run-python benchmark test clean

default: benchmark

build-go:
	go build -o bench-go ./cmd/bench

run-go: build-go
	./bench-go --limit=500000 > go_output.txt

run-python:
	python3 python/bench.py --limit=500000 > python_output.txt

benchmark: run-go run-python
	./run_bench.sh

test: benchmark
	chmod +x tests/validate_output.sh
	./tests/validate_output.sh

clean:
	rm -f bench-go go_output.txt python_output.txt
