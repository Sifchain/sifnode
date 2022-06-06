#!/usr/bin/env bash

set -x

sifnoded tx clp pmtp-params \
  --pmtp_start=31 \
  --pmtp_end=1030 \
  --epochLength=100 \
  --rGov=0.10 \
  --from=$SIF_ACT \
  --keyring-backend=test \
  --fees 100000000000000000rowan \
  --node ${SIFNODE_NODE} \
  --chain-id=$SIFNODE_CHAIN_ID \
  --broadcast-mode=block \
  -y