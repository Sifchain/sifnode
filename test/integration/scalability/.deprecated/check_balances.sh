#!/bin/bash

echo "check balances of ${1} against all chains"

CHAINS_BINARY="sifnoded,${CHAINS_BINARY}"
CHAINS_NODE="${SIF_NODE},${CHAINS_NODE}"
CHAINS_ID="${SIF_CHAINID},${CHAINS_ID}"
CHAINS_BINARY=(${CHAINS_BINARY//,/ })
CHAINS_NODE=(${CHAINS_NODE//,/ })
CHAINS_ID=(${CHAINS_ID//,/ })

for i in ${!CHAINS_BINARY[@]}; do
    echo "${CHAINS_BINARY[$i]}"
    echo "------------------"
    ${CHAINS_BINARY[$i]} \
        q \
        bank \
        balances \
        $(./get_address.sh ${CHAINS_BINARY[$i]} $1) \
        --node ${CHAINS_NODE[$i]} \
        --chain-id=${CHAINS_ID[$i]}
    echo
done