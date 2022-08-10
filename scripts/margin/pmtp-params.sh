#!/usr/bin/env bash

set -x

sifnoded q clp pmtp-params \
  --node ${SIFNODE_NODE} \
  --chain-id $SIFNODE_CHAIN_ID