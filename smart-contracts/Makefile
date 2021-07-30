# Makefile

SHELL:=/bin/bash

setup:
	python3 -m venv venv ; source ./venv/bin/activate ; pip3 install -r requirements.txt
	yarn

# Run slither over entire directory
# `make slither`
slither: setup
	slither . || true

# Simple static analysis in a human-readable report over entire directory
# `make slither-pretty-summary`
slither-pretty-summary: setup
	slither . --print human-summary

# Check for ERC 20|223|777|721|165|1820 conformance
# `make conformance CONTRACT=./contracts/ContractFile.sol CONTRACT_NAME=ContractName`
erc-conformance: setup
	slither-check-erc ${CONTRACT} ${CONTRACT_NAME}

.PHONY: clean clean-smartcontracts clean-node
clean: clean-node clean-smartcontracts

clean-smartcontracts:
	rm -rf build .openzepplin

clean-node:
	rm -rf node_modules

abi:
	yarn peggy:generateAbi