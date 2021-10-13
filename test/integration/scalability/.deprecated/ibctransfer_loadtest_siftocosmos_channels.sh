#!/bin/bash -x

SIF_IBC_ADDRESS=$(sifnoded keys show $1 --keyring-backend test -a)
COSMOS_IBC_ADDRESS=$(gaiad keys show $1 --keyring-backend test -a)
RESPONSE=$(sifnoded q auth account ${SIF_IBC_ADDRESS} --node ${SIF_NODE} --chain-id ${SIF_CHAINID} --output json)
ACCOUNT_NUMBER=$(echo $RESPONSE | jq -r .account_number)
SEQUENCE=$(echo $RESPONSE | jq -r .sequence)

seq=$SEQUENCE
for channel_id in ${SIF_CHANNEL_IDS//,/ }; do
    echo "channel_id ${channel_id} processing"
    sifnoded \
        tx \
        ibc-transfer \
        transfer \
        transfer \
        ${channel_id} \
        ${COSMOS_IBC_ADDRESS} \
        1rowan \
        --from ${SIF_IBC_ADDRESS} \
        --keyring-backend test \
        --fees 100000rowan \
        --chain-id ${COSMOS_CHAINID} \
        --node ${COSMOS_NODE} \
        --broadcast-mode async \
        --packet-timeout-timestamp 0 \
        --offline \
        --sequence $seq \
        --account-number $ACCOUNT_NUMBER \
        -y
    echo "channel_id ${channel_id} done"
    seq=$((seq + 1))
done