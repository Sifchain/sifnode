#!/bin/sh

TOKEN_REGISTRY_ADMIN_ADDRESS="sif1tpypxpppcf5lea47vcvgy09675nllmcucxydvu"

sifnoded tx tokenregistry deregister-all ./$SIFCHAIN_ID/tokenregistry.json \
  --node $SIF_NODE \
  --chain-id $SIFCHAIN_ID \
  --from $TOKEN_REGISTRY_ADMIN_ADDRESS \
  --gas-prices=0.5rowan \
  --gas-adjustment=1.5 \
  --broadcast-mode=block