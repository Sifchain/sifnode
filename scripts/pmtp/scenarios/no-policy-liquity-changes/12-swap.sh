#!/usr/bin/env bash

set -x

sifnoded tx clp swap \
  --from $SIF_ACT \
  --keyring-backend test \
  --sentSymbol rowan \
  --receivedSymbol cusdt \
  --sentAmount 100000000000000000000000 \
  --minReceivingAmount 0 \
  --fees 100000000000000000rowan \
  --node ${SIFNODE_NODE} \
  --chain-id $SIFNODE_CHAIN_ID \
  --broadcast-mode block \
  -y