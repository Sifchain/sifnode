#!/usr/bin/env bash

set -x

ACCOUNT_NUMBER=$(sifnoded q auth account $ADMIN_ADDRESS \
    --node ${SIFNODE_NODE} \
    --chain-id $SIFNODE_CHAIN_ID \
    --output json \
    | jq -r ".account_number")

for i in {0..1}; do
  echo "tx ${i}"
  SEQUENCE=$(sifnoded q auth account $ADMIN_ADDRESS \
    --node ${SIFNODE_NODE} \
    --chain-id $SIFNODE_CHAIN_ID \
    --output json \
    | jq -r ".sequence")
  sifnoded tx clp swap \
    --from=$SIF_ACT \
    --keyring-backend=test \
    --sentSymbol=ibc/27394FB092D2ECCD56123C74F36E4C1F926001CEADA9CA97EA622B25F41E5EB2 \
    --receivedSymbol=rowan \
    --sentAmount=10000000 \
    --minReceivingAmount=0 \
    --fees=100000000000000000rowan \
    --gas=500000 \
    --node=${SIFNODE_NODE} \
    --chain-id=${SIFNODE_CHAIN_ID} \
    --broadcast-mode=block \
    --account-number=${ACCOUNT_NUMBER} \
    --sequence=$SEQUENCE \
    -y
    sleep 1
  done