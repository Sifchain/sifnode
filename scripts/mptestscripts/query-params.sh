#!/usr/bin/env bash

set -x

sifnoded query clp params  \
  --chain-id $SIFNODE_CHAIN_ID \
  --node ${SIFNODE_NODE} \
