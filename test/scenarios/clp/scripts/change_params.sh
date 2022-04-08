#!/bin/sh

sifnoded tx clp reward-params \
--from sif --keyring-backend test \
--lock-period 10 \
--cancel-period 720 \
--default-multiplier 0 \
--chain-id localnet \
--broadcast-mode block \
-y
