#!/usr/bin/env bash

set -x

sifnoded tx margin update-params \
  --force-close-threshold=0.10 \
  --health-gain-factor=1.0 \
  --interest-rate-decrease=0.10 \
  --interest-rate-increase=0.10 \
  --interest-rate-max=3.0 \
  --interest-rate-min=0.005 \
  --leverage-max=2 \
  --epoch-length=1 \
  --from=$SIF_ACT \
  --keyring-backend=test \
  --fees 100000000000000000rowan \
  --gas 500000 \
  --node ${SIFNODE_NODE} \
  --chain-id=$SIFNODE_CHAIN_ID \
  --broadcast-mode=block \
  -y