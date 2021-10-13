#!/bin/bash -x

echo "check balance of ${1}"
sifnoded \
    q \
    bank \
    balances \
    $(./get_sif_address.sh $1) \
    --node ${SIF_NODE} \
    --chain-id=${SIF_CHAINID}