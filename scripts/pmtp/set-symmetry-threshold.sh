#!/usr/bin/env bash

set -x

sifnoded tx clp set-symmetry-threshold \
  --threshold=0.000000005 \
  --from=$SIF_ACT \
  --keyring-backend=test \
  --fees=100000000000000000rowan \
  --gas=500000 \
  --node=${SIFNODE_NODE} \
  --chain-id=$SIFNODE_CHAIN_ID \
  --broadcast-mode=block \
  -y