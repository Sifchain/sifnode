#!/usr/bin/env bash

set -x

sifnoded tx clp cancel-unbond \
  --from $SIF_ACT \
  --keyring-backend test \
  --symbol cusdt \
  --units 1000000000 \
  --fees 100000000000000000rowan \
  --node ${SIFNODE_NODE} \
  --chain-id $SIFNODE_CHAIN_ID \
  --broadcast-mode block \
  -y