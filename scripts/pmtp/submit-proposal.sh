#!/usr/bin/env bash

set -x

sifnoded tx gov submit-proposal \
    param-change proposal.json \
    --from $SIF_ACT --keyring-backend test \
    --node tcp://${SIFNODE_P2P_HOSTNAME}:26657 \
    --chain-id $SIFNODE_CHAIN_ID \
    -y \
    --broadcast-mode block \
    --fees 100000000000000000rowan