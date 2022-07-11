#!/usr/bin/env bash

set -x

sifnoded tx clp reward-params \
  --from $SIF_ACT \
  --cancelPeriod 100 \
  --lockPeriod 10 \
  --keyring-backend test \
  --fees 100000000000000000rowan \
  --chain-id $SIFNODE_CHAIN_ID \
  --node ${SIFNODE_NODE} \
  --broadcast-mode block \
  -y
