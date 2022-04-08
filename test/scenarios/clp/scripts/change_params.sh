#!/bin/sh

sifnoded tx clp reward-params \
--from sif --keyring-backend test \
--lockPeriod 10 \
--cancelPeriod 720 \
--defaultMultiplier 0 \
--chain-id localnet \
--broadcast-mode block \
-y
