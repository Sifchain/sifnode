# Makefile

SHELL:=/bin/bash

include Makefile.smartContracts
include Makefile.artifacts
gofiles=$(artifacts:.json=.go)

artifacts/%.go: artifacts/%.json
	bash scripts/do_abigen.sh $< ../cmd/ebrelayer/contract/generated

all: ${gofiles}

node_modules:
	npm install

.PHONY: clean clean-smartcontracts clean-node
clean: clean-node clean-smartcontracts

clean-smartcontracts:
	rm -rf build .openzepplin

clean-node:
	rm -rf node_modules

abi:
	yarn peggy:generateAbi
