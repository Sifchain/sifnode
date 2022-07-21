#!/usr/bin/env bash

set -x

sifnoded tx margin open \
  --from $SIF_ACT \
  --keyring-backend test \
  --borrow_asset cusdc \
  --collateral_asset rowan \
  --collateral_amount 1000000000000000000000 \
  --position long \
  --leverage 2 \
  --fees 100000000000000000rowan \
  --node ${SIFNODE_NODE} \
  --chain-id $SIFNODE_CHAIN_ID \
  --broadcast-mode block \
  -y