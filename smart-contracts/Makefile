# Makefile

SHELL:=/bin/bash

include Makefile.smartContracts
include Makefile.artifacts
gofiles=$(artifacts:.json=.go)

all: ${artifacts}

install: all

artifacts/%.go: artifacts/%.json
	bash scripts/do_abigen.sh $< ../cmd/ebrelayer/contract/generated

${artifacts}: .hardhat-compile

all: ${gofiles}

.hardhat-compile: node_modules
	npx hardhat compile --force
	touch $@

node_modules: package.json
	npm install

.PHONY: clean clean-smartcontracts clean-node clean-hardhat-artifact 
clean: clean-hardhat-artifact clean-node clean-smartcontracts
	rm -f .hardhat-compile

clean-smartcontracts:
	rm -rf build .openzepplin

clean-node:
	rm -rf node_modules

clean-hardhat-artifact:
	git clean -fdx artifacts 
	git clean -fdx cache

type: 
	npx hardhat typechain

tests: type
	npx hardhat test test/*.js test/*.ts
