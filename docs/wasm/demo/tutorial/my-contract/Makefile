build:
	RUSTFLAGS='-C link-arg=-s' cargo build --release --target wasm32-unknown-unknown
	cp target/wasm32-unknown-unknown/release/my_contract.wasm .
.PHONY: build
