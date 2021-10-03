#!/bin/sh

SIFCHAIN_ID=sifchain-testnet-1 \
  KEYRING_BACKEND=test \
  SIF_NODE=https://rpc-testnet.sifchain.finance:443 ./template/register-all-eth.sh