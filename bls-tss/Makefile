ROOT_DIR := $(dir $(realpath $(lastword $(MAKEFILE_LIST))))

.PHONY: build
build: linux-amd64 linux-arm64 darwin-amd64 darwin-arm64

.PHONY: linux-amd64
linux-amd64:
	cargo build --target=x86_64-unknown-linux-gnu --verbose --release
	cp target/x86_64-unknown-linux-gnu/release/libbls_tss.a lib/linux/amd64

.PHONY: linux-arm64
linux-arm64:
	cross build --target=aarch64-unknown-linux-gnu --verbose --release
	cp target/aarch64-unknown-linux-gnu/release/libbls_tss.a lib/linux/arm64

.PHONY: darwin-amd64
darwin-amd64:
	cargo build --target=x86_64-apple-darwin --verbose --release
	cp target/x86_64-apple-darwin/release/libbls_tss.a lib/darwin/amd64

.PHONY: darwin-arm64
darwin-arm64:
	cargo build --target=aarch64-apple-darwin --verbose --release
	cp target/aarch64-apple-darwin/release/libbls_tss.a lib/darwin/arm64

.PHONY: clean
clean:
	rm -rf target
