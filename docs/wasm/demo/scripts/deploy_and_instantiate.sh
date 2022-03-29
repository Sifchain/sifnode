#!/bin/bash

sifnoded tx wasm store ./contract/swapper.wasm \
--from sif --keyring-backend test \
--gas 1000000000000000000 \
--broadcast-mode block \
--chain-id localnet  \
-y

sifnoded tx wasm instantiate 1 '{}' \
--amount 50000rowan \
--label "swapper" \
--no-admin \
--from sif --keyring-backend test \
--gas 1000000000000000000 \
--broadcast-mode block \
--chain-id localnet \
-y