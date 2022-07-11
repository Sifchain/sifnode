#!/usr/bin/env bash

set -x

sifnoded tx clp swap \
  --from $SIF_ACT \
  --keyring-backend test \
  --sentSymbol cusdc \
  --receivedSymbol rowan \
  --sentAmount 184515000000 \
  --minReceivingAmount 0 \
  --fees 100000000000000000rowan \
  --chain-id $SIFNODE_CHAIN_ID \
  --node ${SIFNODE_NODE} \
  --broadcast-mode block \
  -y
