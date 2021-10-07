#!/bin/sh

SIFCHAIN_ID=sifchain-1 \
  KEYRING_BACKEND=test \
  SIF_NODE=https://rpc.sifchain.finance:443 ./template/register-all.sh