#!/usr/bin/env bash

set -x

sifnoded q clp pool cusdt \
  --chain-id $SIFNODE_CHAIN_ID \
  --node ${SIFNODE_NODE} \

sifnoded q clp pool clink \
  --chain-id $SIFNODE_CHAIN_ID \
  --node ${SIFNODE_NODE} \
