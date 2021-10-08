#!/bin/bash -x

echo "check balance of ${1}"
gaiad \
    q \
    bank \
    balances \
    $(./get_cosmos_address.sh $1) \
    --node ${COSMOS_NODE} \
    --chain-id ${COSMOS_CHAINID}