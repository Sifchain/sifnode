#!/bin/sh

# Remove liquidity 
sifnoded tx clp remove-liquidity \
--from sif --keyring-backend test \
--fees 100000000000000000rowan \
--symbol ceth \
--wBasis 5000 --asymmetry 0 \
--chain-id localnet \
--broadcast-mode block \
-y