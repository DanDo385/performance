.PHONY: default build-go build-ts build-cpp build-rust run-go run-python run-ts run-cpp run-rust benchmark test clean install-deps

default: benchmark

install-deps:
	cd typescript && npm install

build-go:
	go build -o bench-go ./cmd/bench

build-ts: install-deps
	cd typescript && npm run build

build-cpp:
	clang++ -O3 -std=c++17 -pthread -o bench-cpp cpp/src/bench.cpp

build-rust:
	cd rust && cargo build --release && cp target/release/bench ../bench-rust

run-go: build-go
	./bench-go --limit=500000 > go_output.txt

run-python:
	python3 python/bench.py --limit=500000 > python_output.txt

run-ts: build-ts
	cd typescript && node dist/bench.js --limit=500000 > ../ts_output.txt

run-cpp: build-cpp
	./bench-cpp --limit=500000 > cpp_output.txt

run-rust: build-rust
	./bench-rust --limit=500000 > rust_output.txt

benchmark: run-go run-python run-ts run-cpp run-rust
	./run_bench.sh

test: benchmark
	chmod +x tests/validate_output.sh
	./tests/validate_output.sh

clean:
	rm -f bench-go bench-cpp bench-rust go_output.txt python_output.txt ts_output.txt cpp_output.txt rust_output.txt
	rm -rf typescript/dist typescript/node_modules
	rm -rf rust/target
