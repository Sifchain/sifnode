#!/usr/bin/env bash

set -x

sifnoded tx clp pmtp-params \
  --pmtp_start=22811 \
  --pmtp_end=224410 \
  --epochLength=14400 \
  --rGov=0.05 \
  --from=$SIF_ACT \
  --keyring-backend=test \
  --fees 100000000000000000rowan \
  --node ${SIFNODE_NODE} \
  --chain-id=$SIFNODE_CHAIN_ID \
  --broadcast-mode=block \
  -y