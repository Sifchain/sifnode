#!/usr/bin/env bash

set -x

sifnoded tx clp set-lppd-params \
  --path=distribution-period.json \
  --from=$SIF_ACT \
  --keyring-backend=test \
  --fees 100000000000000000rowan \
  --chain-id=$SIFNODE_CHAIN_ID \
  --node ${SIFNODE_NODE} \
  --broadcast-mode=block \
  -y
