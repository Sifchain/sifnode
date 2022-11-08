#!/usr/bin/env bash

set -x

sifnoded tx gov submit-proposal \
    param-change proposal.json \
    --from $SIF_ACT \
    --keyring-backend test \
    --node ${SIFNODE_NODE} \
    --chain-id $SIFNODE_CHAIN_ID \
    --fees 100000000000000000rowan \
    --broadcast-mode block \
    -y