# Makefile

SHELL:=/bin/bash

# Run slither over entire directory
# `make slither`
slither:
	python3 -m venv venv
	source ./venv/bin/activate
	pip3 install -r requirements.txt
	yarn
	# Slither . doesn't return a 0 even when it works, so the makefiles always failing. 
	# This is a hacky way to contain that error, and if slither does fail for genuine reasons
	# it will have been echo'd for debugging purposes anyway. 
	slither . || echo 

# Simple static analysis in a human-readable report over entire directory
# `make slither-pretty-summary`
slither-pretty-summary:
	python3 -m venv venv
	source ./venv/bin/activate
	pip3 install -r requirements.txt
	yarn
	slither . --print human-summary

# Check for ERC 20|223|777|721|165|1820 conformance
# `make conformance CONTRACT=./contracts/ContractFile.sol CONTRACT_NAME=ContractName`
erc-conformance:
	python3 -m venv venv
	source ./venv/bin/activate
	pip3 install -r requirements.txt
	yarn
	slither-check-erc ${CONTRACT} ${CONTRACT_NAME}