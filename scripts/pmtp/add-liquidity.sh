#!/usr/bin/env bash

set -x

sifnoded tx clp add-liquidity \
  --from $SIF_ACT \
  --keyring-backend test \
  --symbol ceth \
  --nativeAmount 0 \
  --externalAmount 455000000000000 \
  --fees 100000000000000000rowan \
  --node ${SIFNODE_NODE} \
  --chain-id $SIFNODE_CHAIN_ID \
  --broadcast-mode block \
  -y