#!/usr/bin/env bash

set -x

sifnoded tx clp set-symmetry-threshold \
  --threshold=115792089237316195423570985008687907853269984665640564039457.584007913129639935\
  --ratio=115792089237316195423570985008687907853269984665640564039457.584007913129639935 \
  --from=$SIF_ACT \
  --keyring-backend=test \
  --fees 100000000000000000rowan \
  --chain-id=$SIFNODE_CHAIN_ID \
  --node ${SIFNODE_NODE} \
  --broadcast-mode=block \
  -y
