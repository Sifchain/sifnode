#!/usr/bin/env bash

set -x

sifnoded tx tokenregistry deregister-all denoms-all.json \
  --chain-id "${SIFNODE_CHAIN_ID}" \
  --node ${SIFNODE_NODE} \
  --from "${ADMIN_ADDRESS}" \
  --keyring-backend test \
  --gas 500000 \
  --gas-prices 0.5rowan \
  -y \
  --broadcast-mode block
