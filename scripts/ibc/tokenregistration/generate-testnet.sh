#!/bin/sh

SIFCHAIN_ID=sifchain-testnet-1 \
  COSMOS_BASE_DENOM=uphoton \
  COSMOS_CHANNEL_ID=channel-11 \
  COSMOS_COUNTERPARTY_CHANNEL_ID=channel-27 \
  COSMOS_CHAIN_ID=cosmoshub-testnet \
  AKASH_CHANNEL_ID=channel-12 \
  AKASH_COUNTERPARTY_CHANNEL_ID=channel-66 \
  AKASH_CHAIN_ID=akash-testnet-6 \
  SENTINEL_CHANNEL_ID=channel-13 \
  SENTINEL_COUNTERPARTY_CHANNEL_ID=channel-39 \
  SENTINEL_CHAIN_ID=sentinelhub-2 ./template/generate-all-ibc.sh
