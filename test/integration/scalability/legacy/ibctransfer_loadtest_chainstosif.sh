#!/bin/bash -x

CHAINS_BINARY=(${CHAINS_BINARY//,/ })
CHAINS_NODE=(${CHAINS_NODE//,/ })
CHAINS_ID=(${CHAINS_ID//,/ })
CHAINS_DENOM=(${CHAINS_DENOM//,/ })
CHAINS_FEES=(${CHAINS_FEES//,/ })
CHAINSTOSIF_CHANNEL_ID=(${CHAINSTOSIF_CHANNEL_ID//,/ })

SIF_IBC_ADDRESS=$(sifnoded keys show $1 --keyring-backend test -a)

for i in ${!CHAINS_BINARY[@]}; do
    CHAIN_IBC_ADDRESS=$(${CHAINS_BINARY[$i]} keys show $1 --keyring-backend test -a)

    RESPONSE=$(${CHAINS_BINARY[$i]} q auth account ${CHAIN_IBC_ADDRESS} --node ${CHAINS_NODE[$i]} --chain-id ${CHAINS_ID[$i]} --output json)
    ACCOUNT_NUMBER=$(echo $RESPONSE | jq -r .account_number)
    SEQUENCE=$(echo $RESPONSE | jq -r .sequence)

    echo "chain id ${CHAINS_ID[$i]} processing"

    seq=$SEQUENCE
    for j in {1..99}; do
        echo "tx ${j} processing"
        ${CHAINS_BINARY[$i]} \
            tx \
            ibc-transfer \
            transfer \
            transfer \
            ${CHAINSTOSIF_CHANNEL_ID[$i]} \
            ${SIF_IBC_ADDRESS} \
            1${CHAINS_DENOM[$i]} \
            --from ${CHAIN_IBC_ADDRESS} \
            --keyring-backend test \
            --fees ${CHAINS_FEES[$i]}${CHAINS_DENOM[$i]} \
            --chain-id ${CHAINS_ID[$i]} \
            --node ${CHAINS_NODE[$i]} \
            --broadcast-mode async \
            --packet-timeout-timestamp 0 \
            --offline \
            --sequence $seq \
            --account-number $ACCOUNT_NUMBER \
            -y
        echo "tx ${j} done"
        seq=$((seq + 1))
    done
done

paplay /usr/share/sounds/sound-icons/hash