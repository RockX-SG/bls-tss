ROOT_DIR := $(dir $(realpath $(lastword $(MAKEFILE_LIST))))

# PHONY means that it doesn't correspond to a file; it always runs the build commands.

.PHONY: build-all
build-all: build-demo

.PHONY: run-all
run-all: run-demo

.PHONY: build-demo
build-demo:
	cd lib && cargo build --release
	cp lib/target/release/libbls_tss.so lib/target/release/libbls_tss.a lib/
	go build -ldflags="-r $(ROOT_DIR)lib" examples/demo.go

.PHONY: run-demo
run-demo: build-demo
	RUST_LOG=trace ./demo

.PHONY: clean
clean:
	rm -rf demo lib/libbls_tss.so lib/libbls_tss.a lib/target
