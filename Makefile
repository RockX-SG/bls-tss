ROOT_DIR := $(dir $(realpath $(lastword $(MAKEFILE_LIST))))

# PHONY means that it doesn't correspond to a file; it always runs the build commands.

.PHONY: build-all
build-all: build-demo

.PHONY: run-all
run-all: run-demo

.PHONY: build-demo
build-demo:
	go build -o demo examples/demo.go

.PHONY: run-demo
run-demo: build-demo
	RUST_LOG=trace ./demo

.PHONY: clean
clean:
	rm -rf demo
