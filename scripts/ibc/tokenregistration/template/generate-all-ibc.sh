#!/bin/sh

# REMEMBER to use right counterparty network denom,
# i.e for BetaNet use MAINNET denom registered on counterparty chain, not denom registered on counterparty TESTNET
# i.e for BetaNet, uatom not uphoton, and for TestNet uphoton not uatom.

# Specify these variables when running - see ./run-testnet.sh
#SIFCHAIN_ID=""

#COSMOS_BASE_DENOM
#COSMOS_CHANNEL_ID="channel-"
#COSMOS_COUNTERPARTY_CHANNEL_ID="channel-"
#COSMOS_CHAIN_ID=""

#AKASH_CHANNEL_ID="channel-"
#AKASH_COUNTERPARTY_CHANNEL_ID="channel-"
#AKASH_CHAIN_ID=""

#SENTINEL_CHANNEL_ID="channel-"
#SENTINEL_COUNTERPARTY_CHANNEL_ID="channel-"
#SENTINEL_CHAIN_ID=""

echo "\n\ngenerating and storing all entries for network $SIFCHAIN_ID"

mkdir -p ./$SIFCHAIN_ID

sifnoded q tokenregistry generate \
	--token_base_denom=$COSMOS_BASE_DENOM \
	--token_ibc_counterparty_chain_id=$COSMOS_CHAIN_ID \
  --token_ibc_channel_id=$COSMOS_CHANNEL_ID \
  --token_ibc_counterparty_channel_id=$COSMOS_COUNTERPARTY_CHANNEL_ID \
	--token_ibc_counterparty_denom="" \
	--token_unit_denom="" \
	--token_decimals=6 \
	--token_display_name="" \
	--token_external_symbol="" \
	--token_permission_clp=true \
	--token_permission_ibc_export=true \
	--token_permission_ibc_import=true | jq > $SIFCHAIN_ID/cosmos.json

echo "\n\ngenerated entry for $COSMOS_CHAIN_ID"

cat $SIFCHAIN_ID/cosmos.json | jq

sifnoded q tokenregistry generate \
	--token_base_denom=uakt \
  --token_ibc_counterparty_chain_id=$AKASH_CHAIN_ID \
  --token_ibc_channel_id=$AKASH_CHANNEL_ID \
	--token_ibc_counterparty_channel_id=$AKASH_COUNTERPARTY_CHANNEL_ID \
	--token_ibc_counterparty_denom="" \
	--token_unit_denom="" \
	--token_decimals=6 \
	--token_display_name="UAKT" \
	--token_external_symbol="uakt" \
	--token_permission_clp=true \
	--token_permission_ibc_export=true \
	--token_permission_ibc_import=true | jq > $SIFCHAIN_ID/akash.json

echo "\n\ngenerated entry for $AKASH_CHAIN_ID"

cat $SIFCHAIN_ID/akash.json | jq

sifnoded q tokenregistry generate \
	--token_base_denom=udvpn \
	--token_ibc_counterparty_chain_id=$SENTINEL_CHAIN_ID \
  --token_ibc_channel_id=$SENTINEL_CHANNEL_ID \
  --token_ibc_counterparty_channel_id=$SENTINEL_COUNTERPARTY_CHANNEL_ID \
	--token_ibc_counterparty_denom="" \
	--token_unit_denom="" \
	--token_decimals=6 \
	--token_display_name="uDVPN" \
	--token_external_symbol="udvpn" \
	--token_permission_clp=true \
	--token_permission_ibc_export=true \
	--token_permission_ibc_import=true | jq > $SIFCHAIN_ID/sentinel.json

echo "\n\ngenerated entry for $SENTINEL_CHAIN_ID"

cat $SIFCHAIN_ID/sentinel.json | jq
