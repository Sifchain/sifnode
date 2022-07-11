#!/usr/bin/env bash

set -x

sifnoded q bank balances $ADMIN_ADDRESS \
   --node ${SIFNODE_NODE} \
   --chain-id $SIFNODE_CHAIN_ID

sifnoded q bank balances sif1l7hypmqk2yc334vc6vmdwzp5sdefygj2ad93p5 \
   --node ${SIFNODE_NODE} \
   --chain-id $SIFNODE_CHAIN_ID
