#!/usr/bin/env bash

set -x

sifnoded tx clp pmtp-params \
  --pmtp_start=101 \
  --pmtp_end=1100 \
  --epochLength=100 \
  --rGov=0.0718 \
  --from=$SIF_ACT \
  --keyring-backend=test \
  --fees 100000000000000000rowan \
  --node ${SIFNODE_NODE} \
  --chain-id=$SIFNODE_CHAIN_ID \
  --broadcast-mode=block \
  -y