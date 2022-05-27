#!/usr/bin/env bash

set -x

sifnoded q gov proposals \
    --node ${SIFNODE_NODE} \
    --chain-id $SIFNODE_CHAIN_ID