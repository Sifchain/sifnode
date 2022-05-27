#!/usr/bin/env bash

set -x

sifnoded q bank balances $ADMIN_ADDRESS \
    --node ${SIFNODE_NODE} \
    --chain-id $SIFNODE_CHAIN_ID