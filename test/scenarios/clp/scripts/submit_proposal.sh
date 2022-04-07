#!/bin/sh

# submit proposal to update clp params
sifnoded tx gov submit-proposal param-change ./scripts/proposal.json \
--from sif --keyring-backend test \
--fees 100000rowan \
--chain-id localnet \
--broadcast-mode block \
-y