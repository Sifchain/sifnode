#!/bin/sh

# Vote yes to accept the proposal
sifnoded tx gov vote 1 yes \
--from sif --keyring-backend test \
--fees 100000rowan \
--chain-id  localnet \
--broadcast-mode block \
-y