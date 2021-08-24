#!/bin/sh

SIFCHAIN_ID=sifchain-testnet-1 \
  COSMOS_BASE_DENOM=uphoton \
  COSMOS_CHANNEL_ID=channel-8 \
  COSMOS_COUNTERPARTY_CHANNEL_ID=channel-20 \
  COSMOS_CHAIN_ID=cosmoshub-testnet \
  AKASH_CHANNEL_ID=channel-6 \
  AKASH_COUNTERPARTY_CHANNEL_ID=channel-60 \
  AKASH_CHAIN_ID=akash-testnet-6 \
  SENTINEL_CHANNEL_ID=channel-7 \
  SENTINEL_COUNTERPARTY_CHANNEL_ID=channel-32 \
  SENTINEL_CHAIN_ID=sentinelhub-2 ./template/generate-all-ibc.sh
