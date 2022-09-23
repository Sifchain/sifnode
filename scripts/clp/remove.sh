#!/bin/sh

set -x

sifnoded tx clp remove-liquidity-units \
  --withdrawUnits 1 \
  --symbol ceth \
  --from $SIF_ACT \
  --keyring-backend test \
  --fees 100000000000000000rowan \
  --node ${SIFNODE_NODE} \
  --chain-id $SIFNODE_CHAIN_ID \
  --broadcast-mode block \
  -y