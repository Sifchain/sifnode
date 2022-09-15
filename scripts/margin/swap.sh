#!/usr/bin/env bash

set -x

sifnoded tx clp swap \
  --from=$SIF_ACT \
  --keyring-backend=test \
  --sentSymbol=cusdc \
  --receivedSymbol=rowan \
  --sentAmount=1000000000000 \
  --minReceivingAmount=0 \
  --fees=100000000000000000rowan \
  --gas=500000 \
  --node=${SIFNODE_NODE} \
  --chain-id=${SIFNODE_CHAIN_ID} \
  --broadcast-mode=block \
  -y