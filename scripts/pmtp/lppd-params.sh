#!/usr/bin/env bash

set -x

sifnoded q clp lppd-params \
  --node ${SIFNODE_NODE} \
  --chain-id $SIFNODE_CHAIN_ID  