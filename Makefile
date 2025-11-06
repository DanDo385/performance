.PHONY: default build-go build-ts run-go run-python run-ts benchmark test clean install-deps

default: benchmark

install-deps:
	cd typescript && npm install

build-go:
	go build -o bench-go ./cmd/bench

build-ts: install-deps
	cd typescript && npm run build

run-go: build-go
	./bench-go --limit=500000 --fib=35 > go_output.txt

run-python:
	python3 python/bench.py --limit=500000 --fib=35 > python_output.txt

run-ts: build-ts
	cd typescript && node dist/bench.js --limit=500000 --fib=35 > ../ts_output.txt

benchmark: run-go run-python run-ts
	./run_bench.sh

test: benchmark
	chmod +x tests/validate_output.sh
	./tests/validate_output.sh

clean:
	rm -f bench-go go_output.txt python_output.txt ts_output.txt
	rm -rf typescript/dist typescript/node_modules
