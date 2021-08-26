#!/bin/sh

SIFCHAIN_ID=sifchain-1 \
  COSMOS_BASE_DENOM=uatom \
  COSMOS_CHANNEL_ID=channel-0 \
  COSMOS_COUNTERPARTY_CHANNEL_ID=channel-192 \
  COSMOS_CHAIN_ID=cosmoshub-4 \
  AKASH_CHANNEL_ID=channel-2   \
  AKASH_COUNTERPARTY_CHANNEL_ID=channel-24 \
  AKASH_CHAIN_ID=akashnet-2 \
  SENTINEL_CHANNEL_ID=channel-1 \
  SENTINEL_COUNTERPARTY_CHANNEL_ID=channel-36 \
  SENTINEL_CHAIN_ID=sentinelhub-2 ./template/generate-all-ibc.sh
