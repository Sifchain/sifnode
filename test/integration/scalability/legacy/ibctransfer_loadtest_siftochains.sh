#!/bin/bash -x

CHAINS_BINARY=(${CHAINS_BINARY//,/ })
CHAINS_NODE=(${CHAINS_NODE//,/ })
CHAINS_ID=(${CHAINS_ID//,/ })
SIFTOCHAINS_CHANNEL_ID=(${SIFTOCHAINS_CHANNEL_ID//,/ })

SIF_IBC_ADDRESS=$(sifnoded keys show $1 --keyring-backend test -a)
RESPONSE=$(sifnoded q auth account ${SIF_IBC_ADDRESS} --node ${SIF_NODE} --chain-id ${SIF_CHAINID} --output json)
ACCOUNT_NUMBER=$(echo $RESPONSE | jq -r .account_number)
SEQUENCE=$(echo $RESPONSE | jq -r .sequence)

seq=$SEQUENCE
for i in ${!CHAINS_BINARY[@]}; do
    CHAIN_IBC_ADDRESS=$(${CHAINS_BINARY[$i]} keys show $1 --keyring-backend test -a)

    echo "chain id ${CHAINS_ID[$i]} processing"

    for j in {1..100}; do
        echo "tx ${j} processing"
        sifnoded \
            tx \
            ibc-transfer \
            transfer \
            transfer \
            ${SIFTOCHAINS_CHANNEL_ID[$i]} \
            ${CHAIN_IBC_ADDRESS} \
            1ibc/C9C7D0BEEA163F1F35F3D916A7EA7099FD39FFBB2AAA8257A34277F0429F52BF \
            --from ${SIF_IBC_ADDRESS} \
            --keyring-backend test \
            --fees 150000rowan \
            --gas 300000 \
            --chain-id ${SIF_CHAINID} \
            --node ${SIF_NODE} \
            --broadcast-mode async \
            --packet-timeout-timestamp 0 \
            --offline \
            --sequence $seq \
            --account-number $ACCOUNT_NUMBER \
            -y
        echo "tx ${j} done"
        seq=$((seq + 1))
    done

    echo "chain id ${CHAINS_ID[$i]} done"
done

paplay /usr/share/sounds/sound-icons/hash