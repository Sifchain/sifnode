#!/usr/bin/env bash

set -x

sifnoded tx bank send \
    $ADMIN_KEY \
    sif1syavy2npfyt9tcncdtsdzf7kny9lh777yqc2nd \
    9299999999750930000rowan \
    --keyring-backend=test \
    --node=${SIFNODE_NODE} \
    --chain-id=$SIFNODE_CHAIN_ID \
    --fees=100000000000000000rowan \
    --broadcast-mode=block \
    -y