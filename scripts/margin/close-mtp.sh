#!/usr/bin/env bash

set -x

sifnoded tx margin close \
  --from $SIF_ACT \
  --id 7 \
  --keyring-backend test \
  --fees 100000000000000000rowan \
  --node ${SIFNODE_NODE} \
  --chain-id $SIFNODE_CHAIN_ID \
  --broadcast-mode block \
  -y