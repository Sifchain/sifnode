#!/bin/sh

SIFCHAIN_ID=sifchain-devnet-1 \
  KEYRING_BACKEND=test \
  SIF_NODE=https://rpc-devnet.sifchain.finance:443 ./template/register-all.sh