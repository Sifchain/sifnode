#!/usr/bin/env bash

set -x

sifnoded tx tokenregistry register denoms/rowan.json \
  --node tcp://"${SIFNODE_P2P_HOSTNAME}":26657 \
  --chain-id "${SIFNODE_CHAIN_ID}" \
  --from "${ADMIN_ADDRESS}" \
  --keyring-backend test \
  --gas 500000 \
  --gas-prices 0.5rowan \
  -y \
  --broadcast-mode block

sifnoded tx tokenregistry register denoms/ceth.json \
  --node tcp://"${SIFNODE_P2P_HOSTNAME}":26657 \
  --chain-id "${SIFNODE_CHAIN_ID}" \
  --from "${ADMIN_ADDRESS}" \
  --keyring-backend test \
  --gas 500000 \
  --gas-prices 0.5rowan \
  -y \
  --broadcast-mode block

sifnoded tx tokenregistry register denoms/cusdc.json \
  --node tcp://"${SIFNODE_P2P_HOSTNAME}":26657 \
  --chain-id "${SIFNODE_CHAIN_ID}" \
  --from "${ADMIN_ADDRESS}" \
  --keyring-backend test \
  --gas 500000 \
  --gas-prices 0.5rowan \
  -y \
  --broadcast-mode block

sifnoded tx tokenregistry register denoms/cusdt.json \
  --node tcp://"${SIFNODE_P2P_HOSTNAME}":26657 \
  --chain-id "${SIFNODE_CHAIN_ID}" \
  --from "${ADMIN_ADDRESS}" \
  --keyring-backend test \
  --gas 500000 \
  --gas-prices 0.5rowan \
  -y \
  --broadcast-mode block

sifnoded tx tokenregistry register denoms/uatom.json \
  --node tcp://"${SIFNODE_P2P_HOSTNAME}":26657 \
  --chain-id "${SIFNODE_CHAIN_ID}" \
  --from "${ADMIN_ADDRESS}" \
  --keyring-backend test \
  --gas 500000 \
  --gas-prices 0.5rowan \
  -y \
  --broadcast-mode block

sifnoded tx tokenregistry register denoms/ujuno.json \
  --node tcp://"${SIFNODE_P2P_HOSTNAME}":26657 \
  --chain-id "${SIFNODE_CHAIN_ID}" \
  --from "${ADMIN_ADDRESS}" \
  --keyring-backend test \
  --gas 500000 \
  --gas-prices 0.5rowan \
  -y \
  --broadcast-mode block

sifnoded tx tokenregistry register denoms/uluna.json \
  --node tcp://"${SIFNODE_P2P_HOSTNAME}":26657 \
  --chain-id "${SIFNODE_CHAIN_ID}" \
  --from "${ADMIN_ADDRESS}" \
  --keyring-backend test \
  --gas 500000 \
  --gas-prices 0.5rowan \
  -y \
  --broadcast-mode block