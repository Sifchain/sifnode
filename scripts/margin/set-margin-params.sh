#!/usr/bin/env bash

set -x

sifnoded tx margin update-params \
  --force-close-threshold=0.000000000001 \
  --health-gain-factor=0.001 \
  --interest-rate-decrease=0.000000000001 \
  --interest-rate-increase=0.000000000001 \
  --interest-rate-max=3.0 \
  --interest-rate-min=0.000000000001 \
  --leverage-max=2 \
  --epoch-length=1 \
  --removal-queue-threshold=0.1 \
  --max-open-positions=10000 \
  --from=$SIF_ACT \
  --keyring-backend=test \
  --fees 100000000000000000rowan \
  --gas 500000 \
  --node ${SIFNODE_NODE} \
  --chain-id=$SIFNODE_CHAIN_ID \
  --broadcast-mode=block \
  -y