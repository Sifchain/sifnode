#!/usr/bin/env bash

set -x

sifnoded q clp pmtp-params \
   --chain-id $SIFNODE_CHAIN_ID \
   --node ${SIFNODE_NODE} \
