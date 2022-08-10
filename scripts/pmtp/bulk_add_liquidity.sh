#!/usr/bin/env bash

set -x

ACCOUNT_NUMBER=$(sifnoded q auth account $ADMIN_ADDRESS \
    --node ${SIFNODE_NODE} \
    --chain-id $SIFNODE_CHAIN_ID \
    --output json \
    | jq -r ".account_number")
SEQUENCE=$(sifnoded q auth account $ADMIN_ADDRESS \
  --node ${SIFNODE_NODE} \
  --chain-id $SIFNODE_CHAIN_ID \
  --output json \
  | jq -r ".sequence")
for i in {0..12244}; do
  echo "tx ${i}"
  sifnoded tx clp add_liquidity \
    --from=$SIF_ACT \
    --keyring-backend=test \
    --externalAmount=${EXTERNAL_AMOUNT} \
    --nativeAmount=${NATIVE_AMOUNT} \
    --symbol=${SYMBOL} \
    --fees=100000000000000000rowan \
    --gas=500000 \
    --node=${SIFNODE_NODE} \
    --chain-id=${SIFNODE_CHAIN_ID} \
    --broadcast-mode=async \
    --account-number=${ACCOUNT_NUMBER} \
    --sequence=$(($SEQUENCE + $i)) \
    -y
  done