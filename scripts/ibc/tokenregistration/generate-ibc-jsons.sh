#!/bin/sh

# sh ./generate-ibc-jsons.sh testnet

. ./envs/$1.sh 

echo "\n\ngenerating and storing all entries for network $SIFCHAIN_ID"

mkdir -p ./$SIFCHAIN_ID

sifnoded q tokenregistry generate -o json \
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

echo "\n\ngenerated entry for COSMOS $COSMOS_CHAIN_ID"

cat $SIFCHAIN_ID/cosmos.json | jq

sifnoded q tokenregistry generate -o json \
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

echo "\n\ngenerated entry for AKASH $AKASH_CHAIN_ID"

cat $SIFCHAIN_ID/akash.json | jq

sifnoded q tokenregistry generate -o json \
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

echo "\n\ngenerated entry for SENTINEL $SENTINEL_CHAIN_ID"

cat $SIFCHAIN_ID/sentinel.json | jq


sifnoded q tokenregistry generate -o json \
	--token_base_denom=$IRIS_BASE_DENOM \
	--token_ibc_counterparty_chain_id=$IRIS_CHAIN_ID \
  --token_ibc_channel_id=$IRIS_CHANNEL_ID \
  --token_ibc_counterparty_channel_id=$IRIS_COUNTERPARTY_CHANNEL_ID \
	--token_ibc_counterparty_denom="" \
	--token_unit_denom="" \
	--token_decimals=6 \
	--token_display_name="" \
	--token_external_symbol="" \
	--token_permission_clp=true \
	--token_permission_ibc_export=true \
	--token_permission_ibc_import=true | jq > $SIFCHAIN_ID/iris.json

echo "\n\ngenerated entry for IRIS $IRIS_CHAIN_ID"

cat $SIFCHAIN_ID/iris.json | jq

sifnoded q tokenregistry generate -o json \
	--token_base_denom=uxprt \
	--token_ibc_counterparty_chain_id=$PERSISTENCE_CHAIN_ID \
  --token_ibc_channel_id=$PERSISTENCE_CHANNEL_ID \
  --token_ibc_counterparty_channel_id=$PERSISTENCE_COUNTERPARTY_CHANNEL_ID \
	--token_ibc_counterparty_denom="" \
	--token_unit_denom="" \
	--token_decimals=6 \
	--token_display_name="uXPRT" \
	--token_external_symbol="uxprt" \
	--token_permission_clp=true \
	--token_permission_ibc_export=true \
	--token_permission_ibc_import=true | jq > $SIFCHAIN_ID/persistence.json

echo "\n\ngenerated entry for PERSISTENCE $PERSISTENCE_CHAIN_ID"

cat $SIFCHAIN_ID/persistence.json | jq

sifnoded q tokenregistry generate -o json \
	--token_base_denom=basecro \
	--token_ibc_counterparty_chain_id=$CRYPTO_ORG_CHAIN_ID \
  --token_ibc_channel_id=$CRYPTO_ORG_CHANNEL_ID \
  --token_ibc_counterparty_channel_id=$CRYPTO_ORG_COUNTERPARTY_CHANNEL_ID \
	--token_ibc_counterparty_denom="" \
	--token_unit_denom="" \
	--token_decimals=8 \
	--token_display_name="CRO" \
	--token_external_symbol="basecro" \
	--token_permission_clp=true \
	--token_permission_ibc_export=true \
	--token_permission_ibc_import=true | jq > $SIFCHAIN_ID/crypto-org.json

echo "\n\ngenerated entry for CRYPTO $CRYPTO_ORG_CHAIN_ID"

cat $SIFCHAIN_ID/crypto-org.json | jq

sifnoded q tokenregistry generate -o json \
	--token_base_denom=uregen \
	--token_ibc_counterparty_chain_id=$REGEN_CHAIN_ID \
  --token_ibc_channel_id=$REGEN_CHANNEL_ID \
  --token_ibc_counterparty_channel_id=$REGEN_COUNTERPARTY_CHANNEL_ID \
	--token_ibc_counterparty_denom="" \
	--token_unit_denom="" \
	--token_decimals=6 \
	--token_display_name="" \
	--token_external_symbol="" \
	--token_permission_clp=true \
	--token_permission_ibc_export=true \
	--token_permission_ibc_import=true | jq > $SIFCHAIN_ID/regen.json

echo "\n\ngenerated entry for REGEN $REGEN_CHAIN_ID"

cat $SIFCHAIN_ID/regen.json | jq

sifnoded q tokenregistry generate -o json \
	--token_base_denom=uosmo \
	--token_ibc_counterparty_chain_id=$OSMOSIS_CHAIN_ID \
  --token_ibc_channel_id=$OSMOSIS_CHANNEL_ID \
  --token_ibc_counterparty_channel_id=$OSMOSIS_COUNTERPARTY_CHANNEL_ID \
	--token_ibc_counterparty_denom="" \
	--token_unit_denom="" \
	--token_decimals=6 \
	--token_display_name="" \
	--token_external_symbol="" \
	--token_permission_clp=true \
	--token_permission_ibc_export=true \
	--token_permission_ibc_import=true | jq > $SIFCHAIN_ID/osmosis.json

echo "\n\ngenerated entry for OSMOSIS $OSMOSIS_CHAIN_ID"

cat $SIFCHAIN_ID/osmosis.json | jq

sifnoded q tokenregistry generate -o json \
	--token_base_denom=ujuno \
	--token_ibc_counterparty_chain_id=$JUNO_CHAIN_ID \
  --token_ibc_channel_id=$JUNO_CHANNEL_ID \
  --token_ibc_counterparty_channel_id=$JUNO_COUNTERPARTY_CHANNEL_ID \
	--token_ibc_counterparty_denom="" \
	--token_unit_denom="" \
	--token_decimals=6 \
	--token_display_name="" \
	--token_external_symbol="" \
	--token_permission_clp=true \
	--token_permission_ibc_export=true \
	--token_permission_ibc_import=true | jq > $SIFCHAIN_ID/juno.json

echo "\n\ngenerated entry for JUNO $JUNO_CHAIN_ID"

cat $SIFCHAIN_ID/juno.json | jq

sifnoded q tokenregistry generate -o json \
	--token_base_denom=uixo \
	--token_ibc_counterparty_chain_id=$IXO_CHAIN_ID \
  --token_ibc_channel_id=$IXO_CHANNEL_ID \
  --token_ibc_counterparty_channel_id=$IXO_COUNTERPARTY_CHANNEL_ID \
	--token_ibc_counterparty_denom="" \
	--token_unit_denom="" \
	--token_decimals=6 \
	--token_display_name="" \
	--token_external_symbol="" \
	--token_permission_clp=true \
	--token_permission_ibc_export=true \
	--token_permission_ibc_import=true | jq > $SIFCHAIN_ID/ixo.json

echo "\n\ngenerated entry for IXO $IXO_CHAIN_ID"

cat $SIFCHAIN_ID/ixo.json | jq

sifnoded q tokenregistry generate -o json \
	--token_base_denom=nanolike \
	--token_ibc_counterparty_chain_id=$LIKECOIN_CHAIN_ID \
  --token_ibc_channel_id=$LIKECOIN_CHANNEL_ID \
  --token_ibc_counterparty_channel_id=$LIKECOIN_COUNTERPARTY_CHANNEL_ID \
	--token_ibc_counterparty_denom="" \
	--token_unit_denom="" \
	--token_decimals=9 \
	--token_display_name="" \
	--token_external_symbol="" \
	--token_permission_clp=true \
	--token_permission_ibc_export=true \
	--token_permission_ibc_import=true | jq > $SIFCHAIN_ID/likecoin.json

echo "\n\ngenerated entry for LIKECOIN $LIKECOIN_CHAIN_ID"

cat $SIFCHAIN_ID/likecoin.json | jq

sifnoded q tokenregistry generate -o json \
	--token_base_denom=ubtsg \
	--token_ibc_counterparty_chain_id=$BITSONG_CHAIN_ID \
  --token_ibc_channel_id=$BITSONG_CHANNEL_ID \
  --token_ibc_counterparty_channel_id=$BITSONG_COUNTERPARTY_CHANNEL_ID \
	--token_ibc_counterparty_denom="" \
	--token_unit_denom="" \
	--token_decimals=6 \
	--token_display_name="" \
	--token_external_symbol="" \
	--token_permission_clp=true \
	--token_permission_ibc_export=true \
	--token_permission_ibc_import=true | jq > $SIFCHAIN_ID/bitsong.json

echo "\n\ngenerated entry for BITSONG $BITSONG_CHAIN_ID"

cat $SIFCHAIN_ID/bitsong.json | jq

sifnoded q tokenregistry generate -o json \
	--token_base_denom=uband \
	--token_ibc_counterparty_chain_id=$BAND_CHAIN_ID \
  --token_ibc_channel_id=$BAND_CHANNEL_ID \
  --token_ibc_counterparty_channel_id=$BAND_COUNTERPARTY_CHANNEL_ID \
	--token_ibc_counterparty_denom="" \
	--token_unit_denom="" \
	--token_decimals=6 \
	--token_display_name="" \
	--token_external_symbol="" \
	--token_permission_clp=true \
	--token_permission_ibc_export=true \
	--token_permission_ibc_import=true | jq > $SIFCHAIN_ID/band.json

echo "\n\ngenerated entry for BAND $BAND_CHAIN_ID"

cat $SIFCHAIN_ID/band.json | jq

sifnoded q tokenregistry generate -o json \
	--token_base_denom=ungm \
	--token_ibc_counterparty_chain_id=$EMONEY_CHAIN_ID \
  --token_ibc_channel_id=$EMONEY_CHANNEL_ID \
  --token_ibc_counterparty_channel_id=$EMONEY_COUNTERPARTY_CHANNEL_ID \
	--token_ibc_counterparty_denom="" \
	--token_unit_denom="" \
	--token_decimals=6 \
	--token_display_name="" \
	--token_external_symbol="" \
	--token_permission_clp=true \
	--token_permission_ibc_export=true \
	--token_permission_ibc_import=true | jq > $SIFCHAIN_ID/emoney.json

echo "\n\ngenerated entry for EMONEY $EMONEY_CHAIN_ID"

cat $SIFCHAIN_ID/emoney.json | jq

sifnoded q tokenregistry generate -o json \
	--token_base_denom=eeur \
	--token_ibc_counterparty_chain_id=$EMONEY_CHAIN_ID \
  --token_ibc_channel_id=$EMONEY_CHANNEL_ID \
  --token_ibc_counterparty_channel_id=$EMONEY_COUNTERPARTY_CHANNEL_ID \
	--token_ibc_counterparty_denom="" \
	--token_unit_denom="" \
	--token_decimals=6 \
	--token_display_name="" \
	--token_external_symbol="" \
	--token_permission_clp=true \
	--token_permission_ibc_export=true \
	--token_permission_ibc_import=true | jq > $SIFCHAIN_ID/emoney-eeur.json

echo "\n\ngenerated entry for EMONEY $EMONEY_CHAIN_ID"

cat $SIFCHAIN_ID/emoney-eeur.json | jq

sifnoded q tokenregistry generate -o json \
	--token_base_denom=uluna \
	--token_ibc_counterparty_chain_id=$TERRA_CHAIN_ID \
  --token_ibc_channel_id=$TERRA_CHANNEL_ID \
  --token_ibc_counterparty_channel_id=$TERRA_COUNTERPARTY_CHANNEL_ID \
	--token_ibc_counterparty_denom="" \
	--token_unit_denom="" \
	--token_decimals=6 \
	--token_display_name="" \
	--token_external_symbol="" \
	--token_permission_clp=true \
	--token_permission_ibc_export=true \
	--token_permission_ibc_import=true | jq > $SIFCHAIN_ID/terra.json

echo "\n\ngenerated entry for TERRA $TERRA_CHAIN_ID"

cat $SIFCHAIN_ID/terra.json | jq

sifnoded q tokenregistry generate -o json \
	--token_base_denom=uusd \
	--token_ibc_counterparty_chain_id=$TERRA_CHAIN_ID \
  --token_ibc_channel_id=$TERRA_CHANNEL_ID \
  --token_ibc_counterparty_channel_id=$TERRA_COUNTERPARTY_CHANNEL_ID \
	--token_ibc_counterparty_denom="" \
	--token_unit_denom="" \
	--token_decimals=6 \
	--token_display_name="" \
	--token_external_symbol="" \
	--token_permission_clp=true \
	--token_permission_ibc_export=true \
	--token_permission_ibc_import=true | jq > $SIFCHAIN_ID/terra-uusd.json

echo "\n\ngenerated entry for TERRA $TERRA_CHAIN_ID"

cat $SIFCHAIN_ID/terra-uusd.json | jq

sifnoded q tokenregistry generate -o json \
	--token_base_denom=uscrt \
	--token_ibc_counterparty_chain_id=$SECRET_CHAIN_ID \
  --token_ibc_channel_id=$SECRET_CHANNEL_ID \
  --token_ibc_counterparty_channel_id=$SECRET_COUNTERPARTY_CHANNEL_ID \
	--token_ibc_counterparty_denom="" \
	--token_unit_denom="" \
	--token_decimals=6 \
	--token_display_name="Secret" \
	--token_external_symbol="" \
	--token_permission_clp=true \
	--token_permission_ibc_export=true \
	--token_permission_ibc_import=true | jq > $SIFCHAIN_ID/scrt.json

echo "\n\ngenerated entry for SECRET $SECRET_CHAIN_ID"

cat $SIFCHAIN_ID/scrt.json | jq

sifnoded q tokenregistry generate -o json \
	--token_base_denom=cdmx \
	--token_ibc_counterparty_chain_id=$COMDEX_CHAIN_ID \
  --token_ibc_channel_id=$COMDEX_CHANNEL_ID \
  --token_ibc_counterparty_channel_id=$COMDEX_COUNTERPARTY_CHANNEL_ID \
	--token_ibc_counterparty_denom="" \
	--token_unit_denom="" \
	--token_decimals=18 \
	--token_display_name="Comdex" \
	--token_external_symbol="" \
	--token_permission_clp=true \
	--token_permission_ibc_export=true \
	--token_permission_ibc_import=true | jq > $SIFCHAIN_ID/cdmx.json

echo "\n\ngenerated entry for $COMDEX_CHAIN_ID"

cat $SIFCHAIN_ID/cmdx.json | jq

sifnoded q tokenregistry generate -o json \
	--token_base_denom=uhuahua \
	--token_ibc_counterparty_chain_id=$HUAHUA_CHAIN_ID \
  --token_ibc_channel_id=$HUAHUA_CHANNEL_ID \
  --token_ibc_counterparty_channel_id=$HUAHUA_COUNTERPARTY_CHANNEL_ID \
	--token_ibc_counterparty_denom="" \
	--token_unit_denom="" \
	--token_decimals=6 \
	--token_display_name="Chihuahua" \
	--token_external_symbol="" \
	--token_permission_clp=true \
	--token_permission_ibc_export=true \
	--token_permission_ibc_import=true | jq > $SIFCHAIN_ID/huahua.json

echo "\n\ngenerated entry for $HUAHUA_CHAIN_ID"

cat $SIFCHAIN_ID/huahua.json | jq

sifnoded q tokenregistry generate -o json \
	--token_base_denom=ustars \
	--token_ibc_counterparty_chain_id=$STARGAZE_CHAIN_ID \
  --token_ibc_channel_id=$STARGAZE_CHANNEL_ID \
  --token_ibc_counterparty_channel_id=$STARGAZE_COUNTERPARTY_CHANNEL_ID \
	--token_ibc_counterparty_denom="" \
	--token_unit_denom="" \
	--token_decimals=6 \
	--token_display_name="Stargaze" \
	--token_external_symbol="" \
	--token_permission_clp=true \
	--token_permission_ibc_export=true \
	--token_permission_ibc_import=true | jq > $SIFCHAIN_ID/ustars.json

echo "\n\ngenerated entry for $STARGAZE_CHAIN_ID"

cat $SIFCHAIN_ID/ustars.json | jq

sifnoded q tokenregistry generate -o json \
	--token_base_denom=ubcna \
	--token_ibc_counterparty_chain_id=$BITCANNA_CHAIN_ID \
  --token_ibc_channel_id=$BITCANNA_CHANNEL_ID \
  --token_ibc_counterparty_channel_id=$BITCANNA_COUNTERPARTY_CHANNEL_ID \
	--token_ibc_counterparty_denom="" \
	--token_unit_denom="" \
	--token_decimals=6 \
	--token_display_name="Bitcanna" \
	--token_external_symbol="" \
	--token_permission_clp=true \
	--token_permission_ibc_export=true \
	--token_permission_ibc_import=true | jq > $SIFCHAIN_ID/ubcna.json

echo "\n\ngenerated entry for $BITCANNA_CHAIN_ID"

cat $SIFCHAIN_ID/ubcna.json | jq

sifnoded q tokenregistry generate -o json \
	--token_base_denom=uiov \
	--token_ibc_counterparty_chain_id=$STARNAME_CHAIN_ID \
  --token_ibc_channel_id=$STARNAME_CHANNEL_ID \
  --token_ibc_counterparty_channel_id=$STARNAME_COUNTERPARTY_CHANNEL_ID \
	--token_ibc_counterparty_denom="" \
	--token_unit_denom="" \
	--token_decimals=6 \
	--token_display_name="Starname" \
	--token_external_symbol="" \
	--token_permission_clp=true \
	--token_permission_ibc_export=true \
	--token_permission_ibc_import=true | jq > $SIFCHAIN_ID/uiov.json

echo "\n\ngenerated entry for $STARNAME_CHAIN_ID"

cat $SIFCHAIN_ID/uiov.json | jq

sifnoded q tokenregistry generate -o json \
	--token_base_denom=aevmos \
	--token_ibc_counterparty_chain_id=$EVMOS_CHAIN_ID \
  --token_ibc_channel_id=$EVMOS_CHANNEL_ID \
  --token_ibc_counterparty_channel_id=$EVMOS_COUNTERPARTY_CHANNEL_ID \
	--token_ibc_counterparty_denom="" \
	--token_unit_denom="" \
	--token_decimals=6 \
	--token_display_name="Evmos" \
	--token_external_symbol="" \
	--token_permission_clp=true \
	--token_permission_ibc_export=true \
	--token_permission_ibc_import=true | jq > $SIFCHAIN_ID/aevmos.json

echo "\n\ngenerated entry for $EVMOS_CHAIN_ID"

cat $SIFCHAIN_ID/aevmos.json | jq
