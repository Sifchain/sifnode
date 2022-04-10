#!/usr/bin/env bash

set -x

sifnoded tx tokenregistry register denoms/rowan.json \
  --node ${SIFNODE_NODE} \
  --chain-id "${SIFNODE_CHAIN_ID}" \
  --from "${ADMIN_ADDRESS}" \
  --keyring-backend test \
  --gas 500000 \
  --gas-prices 0.5rowan \
  -y \
  --broadcast-mode block

sifnoded tx tokenregistry register denoms/ceth.json \
  --node ${SIFNODE_NODE} \
  --chain-id "${SIFNODE_CHAIN_ID}" \
  --from "${ADMIN_ADDRESS}" \
  --keyring-backend test \
  --gas 500000 \
  --gas-prices 0.5rowan \
  -y \
  --broadcast-mode block

sifnoded tx tokenregistry register denoms/cusdc.json \
  --node ${SIFNODE_NODE} \
  --chain-id "${SIFNODE_CHAIN_ID}" \
  --from "${ADMIN_ADDRESS}" \
  --keyring-backend test \
  --gas 500000 \
  --gas-prices 0.5rowan \
  -y \
  --broadcast-mode block

sifnoded tx tokenregistry register denoms/cusdt.json \
  --node ${SIFNODE_NODE} \
  --chain-id "${SIFNODE_CHAIN_ID}" \
  --from "${ADMIN_ADDRESS}" \
  --keyring-backend test \
  --gas 500000 \
  --gas-prices 0.5rowan \
  -y \
  --broadcast-mode block

sifnoded tx tokenregistry register denoms/uatom.json \
  --node ${SIFNODE_NODE} \
  --chain-id "${SIFNODE_CHAIN_ID}" \
  --from "${ADMIN_ADDRESS}" \
  --keyring-backend test \
  --gas 500000 \
  --gas-prices 0.5rowan \
  -y \
  --broadcast-mode block

sifnoded tx tokenregistry register denoms/ujuno.json \
  --node ${SIFNODE_NODE} \
  --chain-id "${SIFNODE_CHAIN_ID}" \
  --from "${ADMIN_ADDRESS}" \
  --keyring-backend test \
  --gas 500000 \
  --gas-prices 0.5rowan \
  -y \
  --broadcast-mode block

sifnoded tx tokenregistry register denoms/uluna.json \
  --node ${SIFNODE_NODE} \
  --chain-id "${SIFNODE_CHAIN_ID}" \
  --from "${ADMIN_ADDRESS}" \
  --keyring-backend test \
  --gas 500000 \
  --gas-prices 0.5rowan \
  -y \
  --broadcast-mode block