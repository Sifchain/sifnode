#!/usr/bin/env bash

set -x

sifnoded q clp pools \
  --node ${SIFNODE_NODE} \
  --chain-id $SIFNODE_CHAIN_ID