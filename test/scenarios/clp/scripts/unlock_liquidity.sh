#!/bin/sh

# Unbond liquidity
sifnoded tx clp unbond-liquidity \
--from sif --keyring-backend test \
--fees 100000000000000000rowan \
--symbol ceth \
--units 1000000000000000000 \
--chain-id localnet \
--broadcast-mode block \
-y