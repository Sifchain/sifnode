#!/usr/bin/env bash

set -x

sifnoded tx bank send \
    $SIF_ACT \
    sif144w8cpva2xkly74xrms8djg69y3mljzplx3fjt \
    9299999999750930000rowan \
    --keyring-backend test \
    --node ${SIFNODE_NODE} \
    --chain-id $SIFNODE_CHAIN_ID \
    --fees 100000000000000000rowan \
    --broadcast-mode block \
    -y