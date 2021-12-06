# Makefile

SHELL:=/bin/bash

include Makefile.smartContracts
include Makefile.artifacts
gofiles=$(artifacts:.json=.go)

all: ${artifacts}

install: all

# x:; echo artifcats: ${artifacts}

artifacts/%.go: artifacts/%.json
	bash scripts/do_abigen.sh $< ../cmd/ebrelayer/contract/generated

${artifacts}: .hardhat-compile

all: ${gofiles}

.hardhat-compile: node_modules
	npx hardhat compile --force
	touch $@

node_modules: package.json
	npm install

.PHONY: clean clean-smartcontracts clean-node
clean: clean-node clean-smartcontracts
	rm -f .hardhat-compile

clean-smartcontracts:
	rm -rf build .openzepplin

clean-node:
	rm -rf node_modules

abi:
	yarn peggy:generateAbi
