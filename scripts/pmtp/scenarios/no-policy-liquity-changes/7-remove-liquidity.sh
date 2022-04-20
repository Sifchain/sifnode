#!/usr/bin/env bash

set -x

sifnoded tx clp remove-liquidity \
  --from $SIF_ACT \
  --keyring-backend test \
  --symbol cusdt \
  --asymmetry 10000 \
  --wBasis 295 \
  --fees 100000000000000000rowan \
  --node ${SIFNODE_NODE} \
  --chain-id $SIFNODE_CHAIN_ID \
  --broadcast-mode block \
  -y