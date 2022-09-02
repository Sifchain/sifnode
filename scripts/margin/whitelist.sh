#!/usr/bin/env bash

set -x

sifnoded tx margin whitelist $ADMIN_ADDRESS \
  --from $ADMIN_KEY \
  --keyring-backend test \
  --fees 100000000000000000rowan \
  --node ${SIFNODE_NODE} \
  --chain-id $SIFNODE_CHAIN_ID \
  --broadcast-mode block \
  -y