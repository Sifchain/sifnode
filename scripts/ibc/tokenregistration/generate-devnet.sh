#!/bin/sh

SIFCHAIN_ID=sifchain-devnet-1 \
  COSMOS_BASE_DENOM=uphoton \
  COSMOS_CHANNEL_ID=channel-114 \
  COSMOS_COUNTERPARTY_CHANNEL_ID=channel-26 \
  COSMOS_CHAIN_ID=cosmoshub-testnet \
  AKASH_CHANNEL_ID=channel-110 \
  AKASH_COUNTERPARTY_CHANNEL_ID=channel-channel-63 \
  AKASH_CHAIN_ID=akash-testnet-6 \
  SENTINEL_CHANNEL_ID=channel-111 \
  SENTINEL_COUNTERPARTY_CHANNEL_ID=CHANNEL-35 \
  SENTINEL_CHAIN_ID=sentinelhub-2 ./template/generate-all-ibc.sh
