#!/usr/bin/env bash

set -x

sifnoded tx margin dewhitelist sif1syavy2npfyt9tcncdtsdzf7kny9lh777yqc2nd \
  --from $SIF_ACT \
  --keyring-backend test \
  --fees 100000000000000000rowan \
  --node ${SIFNODE_NODE} \
  --chain-id $SIFNODE_CHAIN_ID \
  --broadcast-mode block \
  -y