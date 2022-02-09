build:
	RUSTFLAGS='-C link-arg=-s' cargo wasm
	cp target/wasm32-unknown-unknown/release/swapper.wasm .
.PHONY: build
