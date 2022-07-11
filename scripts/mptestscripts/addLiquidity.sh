#!/usr/bin/env bash

set -x

sifnoded tx clp add-liquidity \
  --from akasha \
  --keyring-backend test \
  --symbol cusdt \
  --nativeAmount 3084709842431347056190224 \
  --externalAmount 21065550345 \
  --fees 100000000000000000rowan \
  --chain-id $SIFNODE_CHAIN_ID \
  --node ${SIFNODE_NODE} \
  --broadcast-mode block \
  -y

sifnoded tx clp add-liquidity \
  --from akasha \
  --keyring-backend test \
  --symbol cusdt \
  --nativeAmount 47320899675817318625027 \
  --externalAmount 4840560663984643888 \
  --fees 100000000000000000rowan \
  --chain-id $SIFNODE_CHAIN_ID \
  --node ${SIFNODE_NODE} \
  --broadcast-mode block \
  -y
