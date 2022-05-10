#!/usr/bin/env bash

set -x

sifnoded tx clp add-liquidity \
  --from $SIF_ACT \
  --keyring-backend test \
  --symbol cusdt \
  --nativeAmount 1000000000000000000000000 \
  --externalAmount 25378853317 \
  --fees 100000000000000000rowan \
  --node ${SIFNODE_NODE} \
  --chain-id $SIFNODE_CHAIN_ID \
  --broadcast-mode block \
  -y