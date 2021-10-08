#!/bin/sh

# Sifchain network id
# SIFCHAIN_ID=""
# Sifchain node uri
# SIF_NODE=""
# Sifchain token registry address
TOKEN_REGISTRY_ADMIN_ADDRESS="sif1tpypxpppcf5lea47vcvgy09675nllmcucxydvu"
# Admin's keyring backend with token registry address
#KEYRING_BACKEND=""

# CETH

sifnoded tx tokenregistry register $SIFCHAIN_ID/xeth.json \
  --node $SIF_NODE \
  --chain-id $SIFCHAIN_ID \
  --from $TOKEN_REGISTRY_ADMIN_ADDRESS \
  --keyring-backend $KEYRING_BACKEND \
  --gas-prices=0.5rowan \
  --gas-adjustment=1.5 \
  --broadcast-mode=block

# register conversion denom before setting the link here.
sifnoded tx tokenregistry register $SIFCHAIN_ID/ceth.json \
  --node $SIF_NODE \
  --chain-id $SIFCHAIN_ID \
  --from $TOKEN_REGISTRY_ADMIN_ADDRESS \
  --keyring-backend $KEYRING_BACKEND \
  --gas-prices=0.5rowan \
  --gas-adjustment=1.5 \
  --broadcast-mode=block

# ROWAN

sifnoded tx tokenregistry register $SIFCHAIN_ID/xrowan.json \
  --node $SIF_NODE \
  --chain-id $SIFCHAIN_ID \
  --from $TOKEN_REGISTRY_ADMIN_ADDRESS \
  --keyring-backend $KEYRING_BACKEND \
  --gas-prices=0.5rowan \
  --gas-adjustment=1.5 \
  --broadcast-mode=block

sifnoded tx tokenregistry register $SIFCHAIN_ID/rowan.json \
  --node $SIF_NODE \
  --chain-id $SIFCHAIN_ID \
  --from $TOKEN_REGISTRY_ADMIN_ADDRESS \
  --keyring-backend $KEYRING_BACKEND \
  --gas-prices=0.5rowan \
  --gas-adjustment=1.5 \
  --broadcast-mode=block

# CUSDC

sifnoded tx tokenregistry register $SIFCHAIN_ID/cusdc.json \
  --node $SIF_NODE \
  --chain-id $SIFCHAIN_ID \
  --from $TOKEN_REGISTRY_ADMIN_ADDRESS \
  --keyring-backend $KEYRING_BACKEND \
  --gas-prices=0.5rowan \
  --gas-adjustment=1.5 \
  --broadcast-mode=block

# CUSDT

sifnoded tx tokenregistry register $SIFCHAIN_ID/cusdt.json \
  --node $SIF_NODE \
  --chain-id $SIFCHAIN_ID \
  --from $TOKEN_REGISTRY_ADMIN_ADDRESS \
  --keyring-backend $KEYRING_BACKEND \
  --gas-prices=0.5rowan \
  --gas-adjustment=1.5 \
  --broadcast-mode=block
