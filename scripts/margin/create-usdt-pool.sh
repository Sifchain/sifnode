#!/usr/bin/env bash

set -x

sifnoded tx clp create-pool \
  --from $SIF_ACT \
  --keyring-backend test \
  --symbol cusdt \
  --nativeAmount 1550459183129248235861408 \
  --externalAmount 174248776094 \
  --fees 100000000000000000rowan \
  --node ${SIFNODE_NODE} \
  --chain-id $SIFNODE_CHAIN_ID \
  --broadcast-mode block \
  -y