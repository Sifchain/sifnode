#!/usr/bin/env bash

set -x

sifnoded q clp pool cusdt \
  --node ${SIFNODE_NODE} \
  --chain-id $SIFNODE_CHAIN_ID