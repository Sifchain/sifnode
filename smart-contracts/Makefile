# Makefile

SHELL := /bin/bash

# Run slither over entire directory
# `make slither`
slither:
	source ./venv/bin/activate
	pip install -r requirements.txt
	yarn
	slither .

# Simple static analysis in a human-readable report over entire directory
# `make slither-pretty-summary`
slither-pretty-summary:
	source ./venv/bin/activate
	pip install -r requirements.txt
	yarn
	slither . --print human-summary

# Check for ERC 20|223|777|721|165|1820 conformance
# `make conformance CONTRACT=./contracts/ContractFile.sol CONTRACT_NAME=ContractName`
conformance:
	source ./venv/bin/activate
	pip install -r requirements.txt
	yarn
	slither-check-erc ${CONTRACT} ${CONTRACT_NAME}
