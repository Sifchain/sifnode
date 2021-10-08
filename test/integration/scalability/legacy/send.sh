#!/bin/bash -x

echo "send ${4} from ${2} to ${3}"
${1} \
    tx \
    bank \
    send \
    $(./get_address.sh ${1} $2) \
    $(./get_address.sh ${1} $3) \
    ${4} \
    --keyring-backend test \
    --node ${SIF_NODE} \
    --chain-id ${SIF_CHAINID} \
    -y
    # --fees 100000rowan \