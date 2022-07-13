#!/usr/bin/env bash

set -x

sifnoded q tokenregistry entries \
    --node ${SIFNODE_NODE} \
    --chain-id $SIFNODE_CHAIN_ID | jq