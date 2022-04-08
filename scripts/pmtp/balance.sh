#!/usr/bin/env bash

set -x

sifnoded q bank balances $ADMIN_ADDRESS \
    --node tcp://${SIFNODE_P2P_HOSTNAME}:26657 \
    --chain-id $SIFNODE_CHAIN_ID