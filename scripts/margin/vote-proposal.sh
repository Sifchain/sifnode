#!/usr/bin/env bash

set -x

sifnoded tx gov vote 2 yes \
  --from ${SIF_ACT} \
  --keyring-backend test \
  --chain-id="${SIFNODE_CHAIN_ID}" \
  --node="${SIFNODE_NODE}" \
  --fees=100000000000000000rowan \
  --broadcast-mode=block \
  -y