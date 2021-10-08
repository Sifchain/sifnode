#!/bin/bash -x

SIF_ADDRESS=$(sifnoded keys show $1 --keyring-backend test -a)
RESPONSE=$(sifnoded q auth account ${SIF_ADDRESS} --node ${SIF_NODE} --chain-id ${SIF_CHAINID} --output json)
ACCOUNT_NUMBER=$(echo $RESPONSE | jq -r .account_number)
SEQUENCE=$(echo $RESPONSE | jq -r .sequence)

seq=$SEQUENCE
for j in {1..1}; do
    echo "tx ${j} processing"
    sifnoded \
        tx \
        ethbridge \
        lock \
        ${SIF_ADDRESS} \
        ${ETH_ADDRESS} \
        1 \
        ibc/C782C1DE5F380BC8A5B7D490684894B439D31847A004B271D7B7BA07751E582A \
        40000000000000000 \
        --from ${SIF_ADDRESS} \
        --keyring-backend test \
        --fees 100000rowan \
        --gas 300000 \
        --chain-id ${SIF_CHAINID} \
        --node ${SIF_NODE} \
        --ethereum-chain-id ${ETH_CHAINID} \
        --broadcast-mode block \
        --timeout-height 0 \
        --offline \
        --sequence $seq \
        --account-number $ACCOUNT_NUMBER \
        -y
    echo "tx ${j} done"
    seq=$((seq + 1))
done

paplay /usr/share/sounds/sound-icons/hash