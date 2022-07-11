#!/usr/bin/env bash

set -x

sifnoded q clp lppd-params \
  --chain-id $SIFNODE_CHAIN_ID \
  --node ${SIFNODE_NODE} \
