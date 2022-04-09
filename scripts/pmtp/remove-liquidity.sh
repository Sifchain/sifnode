#!/usr/bin/env bash

set -x

sifnoded tx clp remove-liquidity \
  --from $SIF_ACT \
  --keyring-backend test \
  --symbol cusdt \
  --asymmetry 0 \
  --wBasis 5000 \
  --fees 100000000000000000rowan \
  --node tcp://${SIFNODE_P2P_HOSTNAME}:26657 \
  --chain-id $SIFNODE_CHAIN_ID \
  --broadcast-mode block \
  -y