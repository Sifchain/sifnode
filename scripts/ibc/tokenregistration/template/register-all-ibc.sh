#!/bin/sh

# Sifchain network id
# SIFCHAIN_ID=""
# Sifchain node uri
# SIF_NODE=""
# Sifchain token registry address
TOKEN_REGISTRY_ADMIN_ADDRESS="sif1tpypxpppcf5lea47vcvgy09675nllmcucxydvu"
# Admin's keyring backend with token registry address
# KEYRING_BACKEND=""

# COSMOS HUB
sifnoded tx tokenregistry register ./$SIFCHAIN_ID/cosmos.json \
  --node $SIF_NODE \
  --chain-id $SIFCHAIN_ID \
  --from $TOKEN_REGISTRY_ADMIN_ADDRESS \
  --keyring-backend $KEYRING_BACKEND \
  --gas-prices=0.5rowan \
  --gas-adjustment=1.5 \
  --broadcast-mode=block

# AKASH
sifnoded tx tokenregistry register ./$SIFCHAIN_ID/akash.json \
  --node $SIF_NODE \
  --chain-id $SIFCHAIN_ID \
  --from $TOKEN_REGISTRY_ADMIN_ADDRESS \
  --keyring-backend $KEYRING_BACKEND \
  --gas-prices=0.5rowan \
  --gas-adjustment=1.5 \
  --broadcast-mode=block

# SENTINEL
sifnoded tx tokenregistry register ./$SIFCHAIN_ID/sentinel.json \
  --node $SIF_NODE \
  --chain-id $SIFCHAIN_ID \
  --from $TOKEN_REGISTRY_ADMIN_ADDRESS \
  --keyring-backend $KEYRING_BACKEND \
  --gas-prices=0.5rowan \
  --gas-adjustment=1.5 \
  --broadcast-mode=block

# IRIS
sifnoded tx tokenregistry register ./$SIFCHAIN_ID/iris.json \
  --node $SIF_NODE \
  --chain-id $SIFCHAIN_ID \
  --from $TOKEN_REGISTRY_ADMIN_ADDRESS \
  --keyring-backend $KEYRING_BACKEND \
  --gas-prices=0.5rowan \
  --gas-adjustment=1.5 \
  --broadcast-mode=block

# PERSISTENCE
sifnoded tx tokenregistry register ./$SIFCHAIN_ID/persistence.json \
  --node $SIF_NODE \
  --chain-id $SIFCHAIN_ID \
  --from $TOKEN_REGISTRY_ADMIN_ADDRESS \
  --keyring-backend $KEYRING_BACKEND \
  --gas-prices=0.5rowan \
  --gas-adjustment=1.5 \
  --broadcast-mode=block

# CRYPTO ORG
sifnoded tx tokenregistry register ./$SIFCHAIN_ID/crypto-org.json \
  --node $SIF_NODE \
  --chain-id $SIFCHAIN_ID \
  --from $TOKEN_REGISTRY_ADMIN_ADDRESS \
  --keyring-backend $KEYRING_BACKEND \
  --gas-prices=0.5rowan \
  --gas-adjustment=1.5 \
  --broadcast-mode=block

# REGEN
sifnoded tx tokenregistry register ./$SIFCHAIN_ID/regen.json \
  --node $SIF_NODE \
  --chain-id $SIFCHAIN_ID \
  --from $TOKEN_REGISTRY_ADMIN_ADDRESS \
  --keyring-backend $KEYRING_BACKEND \
  --gas-prices=0.5rowan \
  --gas-adjustment=1.5 \
  --broadcast-mode=block

# TERRA
sifnoded tx tokenregistry register ./$SIFCHAIN_ID/terra.json \
  --node $SIF_NODE \
  --chain-id $SIFCHAIN_ID \
  --from $TOKEN_REGISTRY_ADMIN_ADDRESS \
  --keyring-backend $KEYRING_BACKEND \
  --gas-prices=0.5rowan \
  --gas-adjustment=1.5 \
  --broadcast-mode=block


# OSMOSIS
sifnoded tx tokenregistry register ./$SIFCHAIN_ID/osmosis.json \
  --node $SIF_NODE \
  --chain-id $SIFCHAIN_ID \
  --from $TOKEN_REGISTRY_ADMIN_ADDRESS \
  --keyring-backend $KEYRING_BACKEND \
  --gas-prices=0.5rowan \
  --gas-adjustment=1.5 \
  --broadcast-mode=block