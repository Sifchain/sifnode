#!/bin/bash -x

SIF_IBC_ADDRESS=$(sifnoded keys show $1 --keyring-backend test -a)
COSMOS_IBC_ADDRESS=$(gaiad keys show $1 --keyring-backend test -a)
NODES=(${COSMOS_NODE})
RESPONSE=$(gaiad q auth account ${COSMOS_IBC_ADDRESS} --node ${COSMOS_NODE} --chain-id ${COSMOS_CHAINID} --output json)
ACCOUNT_NUMBER=$(echo $RESPONSE | jq -r .account_number)
SEQUENCE=$(echo $RESPONSE | jq -r .sequence)

seq=$SEQUENCE
for i in {1..100}; do
    echo "tx ${i} processing"
    gaiad \
        tx \
        ibc-transfer \
        transfer \
        transfer \
        ${COSMOSTOSIF_CHANNEL_ID} \
        ${SIF_IBC_ADDRESS} \
        1uphoton \
        --from ${COSMOS_IBC_ADDRESS} \
        --keyring-backend test \
        --chain-id ${COSMOS_CHAINID} \
        --node ${NODES[$(($i % 1))]} \
        --broadcast-mode async \
        --packet-timeout-timestamp 0 \
        --offline \
        --sequence $seq \
        --account-number $ACCOUNT_NUMBER \
        -y
    echo "tx ${i} done"
    seq=$((seq + 1))
done

paplay /usr/share/sounds/sound-icons/hash