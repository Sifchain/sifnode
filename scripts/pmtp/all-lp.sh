#!/usr/bin/env bash

set -x

sifnoded q clp all-lp \
  --node ${SIFNODE_NODE} \
  --chain-id $SIFNODE_CHAIN_ID