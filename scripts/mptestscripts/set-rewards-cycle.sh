#!/usr/bin/env bash

set -x

sifnoded tx clp reward-period \
  --from $SIF_ACT \
  --path rewards.json \
  --keyring-backend test \
  --fees 100000000000000000rowan \
  --chain-id $SIFNODE_CHAIN_ID \
  --node ${SIFNODE_NODE} \
  --broadcast-mode block \
  -y
