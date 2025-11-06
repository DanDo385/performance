.PHONY: default build-go build-ts build-c build-cpp build-rust build-cpu-monitor build-memory-report run-go run-python run-ts run-c run-cpp run-rust benchmark test clean install-deps cpu memory

default: benchmark

install-deps:
	cd typescript && npm install

build-go:
	go build -o bench-go ./cmd/bench

build-ts: install-deps
	cd typescript && npm run build

build-c:
	clang -O3 -pthread -o bench-c c/src/bench.c

build-cpp:
	clang++ -O3 -std=c++17 -pthread -o bench-cpp cpp/src/bench.cpp

build-rust:
	cd rust && cargo build --release && cp target/release/bench ../bench-rust

build-cpu-monitor:
	go build -o cpu-monitor ./cmd/cpu-monitor

build-memory-report:
	go build -o memory-report ./cmd/memory-report

run-go: build-go
	./bench-go --limit=500000 > go_output.txt

run-python:
	python3 python/bench.py --limit=500000 > python_output.txt

run-ts: build-ts
	cd typescript && node dist/bench.js --limit=500000 > ../ts_output.txt

run-c: build-c
	./bench-c --limit=500000 > c_output.txt

run-cpp: build-cpp
	./bench-cpp --limit=500000 > cpp_output.txt

run-rust: build-rust
	./bench-rust --limit=500000 > rust_output.txt

benchmark: run-go run-python run-ts run-c run-cpp run-rust
	./run_bench.sh

test: benchmark
	chmod +x tests/validate_output.sh
	./tests/validate_output.sh

cpu: build-cpu-monitor build-go
	./cpu-monitor

memory: build-memory-report
	./memory-report

clean:
	rm -f bench-go bench-c bench-cpp bench-rust cpu-monitor memory-report go_output.txt python_output.txt ts_output.txt c_output.txt cpp_output.txt rust_output.txt
	rm -rf typescript/dist typescript/node_modules
	rm -rf rust/target
