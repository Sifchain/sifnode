#!/usr/bin/env bash

set -x

sifnoded tx clp set-symmetry-threshold \
  --threshold=1000000000000000000000000000 \
  --ratio=1000000000000000000000000000 \
  --from=$ADMIN_KEY \
  --keyring-backend=test \
  --fees 100000000000000000rowan \
  --gas 500000 \
  --node ${SIFNODE_NODE} \
  --chain-id=$SIFNODE_CHAIN_ID \
  --broadcast-mode=block \
  -y