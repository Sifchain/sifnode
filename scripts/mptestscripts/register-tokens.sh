#!/usr/bin/env bash

set -x

sifnoded tx tokenregistry register denom-ceth.json \
  --chain-id "${SIFNODE_CHAIN_ID}" \
  --node ${SIFNODE_NODE} \
  --from "${ADMIN_ADDRESS}" \
  --keyring-backend test \
  --gas 500000 \
  --gas-prices 0.5rowan \
  -y \
  --broadcast-mode block

sifnoded tx tokenregistry register denom-rowan.json \
  --chain-id "${SIFNODE_CHAIN_ID}" \
  --node ${SIFNODE_NODE} \
  --from "${ADMIN_ADDRESS}" \
  --keyring-backend test \
  --gas 500000 \
  --gas-prices 0.5rowan \
  -y \
  --broadcast-mode block

sifnoded tx tokenregistry register denom-cusdc.json \
  --chain-id "${SIFNODE_CHAIN_ID}" \
  --node ${SIFNODE_NODE} \
  --from "${ADMIN_ADDRESS}" \
  --keyring-backend test \
  --gas 500000 \
  --gas-prices 0.5rowan \
  -y \
  --broadcast-mode block

sifnoded tx tokenregistry register denom-cusdt.json \
  --chain-id "${SIFNODE_CHAIN_ID}" \
  --node ${SIFNODE_NODE} \
  --from "${ADMIN_ADDRESS}" \
  --keyring-backend test \
  --gas 500000 \
  --gas-prices 0.5rowan \
  -y \
  --broadcast-mode block

sifnoded tx tokenregistry register denom-clink.json \
  --chain-id "${SIFNODE_CHAIN_ID}" \
  --node ${SIFNODE_NODE} \
  --from "${ADMIN_ADDRESS}" \
  --keyring-backend test \
  --gas 500000 \
  --gas-prices 0.5rowan \
  -y \
  --broadcast-mode block
