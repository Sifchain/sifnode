#!/usr/bin/env bash

set -x

sifnoded tx margin update-pools ./pools.json \
  --closed-pools ./closed-pools.json \
  --from=$SIF_ACT \
  --keyring-backend=test \
  --fees 100000000000000000rowan \
  --gas 500000 \
  --node ${SIFNODE_NODE} \
  --chain-id=$SIFNODE_CHAIN_ID \
  --broadcast-mode=block \
  -y