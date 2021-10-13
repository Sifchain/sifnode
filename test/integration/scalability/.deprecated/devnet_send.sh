#!/bin/bash -x

echo "send ${3} from ${1} to ${2}"
sifnoded \
    tx \
    bank \
    send \
    $(sifnoded keys show $1 --keyring-backend test -a) \
    $(sifnoded keys show $2 --keyring-backend test -a) \
    $3 \
    --keyring-backend test \
    --node http://rpc-devnet-042-ibc.sifchain.finance:80 \
    --chain-id sifchain-devnet-042-ibc \
    --fees 100000rowan \
    -y