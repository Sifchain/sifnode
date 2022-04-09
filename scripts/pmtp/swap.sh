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
  --node tcp://${SIFNODE_P2P_HOSTNAME}:26657 \
  --chain-id $SIFNODE_CHAIN_ID \
  --broadcast-mode block \
  -y