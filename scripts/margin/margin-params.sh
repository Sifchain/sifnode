#!/usr/bin/env bash

set -x

sifnoded q margin params \
  --node ${SIFNODE_NODE} \
  --chain-id $SIFNODE_CHAIN_ID