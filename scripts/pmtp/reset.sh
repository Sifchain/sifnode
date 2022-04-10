#!/usr/bin/env bash

set -x

sifnoded tx clp pmtp-rates \
  --blockRate=0.00 \
  --runningRate=0.00 \
  --from=$SIF_ACT \
  --keyring-backend=test \
  --fees 100000000000000000rowan \
  --node ${SIFNODE_NODE} \
  --chain-id=$SIFNODE_CHAIN_ID \
  --broadcast-mode=block \
  -y