#!/usr/bin/env bash

set -x

sifnoded q margin \
  positions-for-address $ADMIN_ADDRESS \
  --node ${SIFNODE_NODE} \
  --chain-id $SIFNODE_CHAIN_ID