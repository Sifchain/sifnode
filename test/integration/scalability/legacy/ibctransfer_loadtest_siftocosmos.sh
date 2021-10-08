#!/bin/bash -x

SIF_IBC_ADDRESS=$(sifnoded keys show $1 --keyring-backend test -a)
COSMOS_IBC_ADDRESS=$(gaiad keys show $1 --keyring-backend test -a)
NODES=(${SIF_NODE})
RESPONSE=$(sifnoded q auth account $SIF_IBC_ADDRESS --node ${SIF_NODE} --chain-id ${SIF_CHAINID} --output json)
ACCOUNT_NUMBER=$(echo $RESPONSE | jq -r .account_number)
SEQUENCE=$(echo $RESPONSE | jq -r .sequence)

seq=$SEQUENCE
for i in {1..1}; do
    echo "tx ${i} processing"
    sifnoded \
        tx \
        ibc-transfer \
        transfer \
        transfer \
        ${SIFTOCOSMOS_CHANNEL_ID} \
        ${COSMOS_IBC_ADDRESS} \
        1rowan \
        --from ${SIF_IBC_ADDRESS} \
        --keyring-backend test \
        --chain-id ${SIF_CHAINID} \
        --node ${NODES[$(($i % 1))]} \
        --fees 100000rowan \
        --broadcast-mode block \
        --packet-timeout-timestamp 0 \
        --offline \
        --sequence $seq \
        --account-number $ACCOUNT_NUMBER \
        -y
    echo "tx ${i} done"
    seq=$((seq + 1))
done

paplay /usr/share/sounds/sound-icons/hash