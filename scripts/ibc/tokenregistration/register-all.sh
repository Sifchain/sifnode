#!/bin/sh
. ./envs/$1.sh 

# sh ./register-all.sh testnet


TOKEN_REGISTRY_ADMIN_ADDRESS="sif1tpypxpppcf5lea47vcvgy09675nllmcucxydvu"

sifnoded tx tokenregistry register-all ./$SIFCHAIN_ID/tokenregistry.json \
  --node $SIF_NODE \
  --chain-id $SIFCHAIN_ID \
  --from $TOKEN_REGISTRY_ADMIN_ADDRESS \
  --keyring-backend $KEYRING_BACKEND \
  --gas=500000 \
  --gas-prices=0.5rowan \
  -y