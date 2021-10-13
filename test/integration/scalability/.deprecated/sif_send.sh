#!/bin/bash -x

echo "send ${3} from ${1} to ${2}"
sifnoded \
    tx \
    bank \
    send \
    $(./get_address.sh sifnoded $1) \
    $(./get_address.sh sifnoded $2) \
    $3 \
    --keyring-backend test \
    --node ${SIF_NODE} \
    --chain-id ${SIF_CHAINID} \
    --fees 100000rowan \
    -y