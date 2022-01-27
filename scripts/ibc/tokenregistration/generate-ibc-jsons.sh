#!/bin/sh

# sh ./generate-ibc-jsons.sh testnet

. ./envs/$1.sh 

echo "\n\ngenerating and storing all entries for network $SIFCHAIN_ID"

mkdir -p ./$SIFCHAIN_ID

OSMOSIS_TOKEN_DECIMALS=6
AKASH_TOKEN_DECIMALS=6
SENTINEL_TOKEN_DECIMALS=6
IRIS_TOKEN_DECIMALS=6
PERSISTENCE_TOKEN_DECIMALS=6
CRYPTO_ORG_TOKEN_DECIMALS=6
REGEN_TOKEN_DECIMALS=6
JUNO_TOKEN_DECIMALS=6
IXO_TOKEN_DECIMALS=6
LIKECOIN_TOKEN_DECIMALS=9
BITSONG_TOKEN_DECIMALS=6
BAND_TOKEN_DECIMALS=6
EMONEY_TOKEN_DECIMALS=6
EMONEY_EUR_TOKEN_DECIMALS=6
TERRA_TOKEN_DECIMALS=6
TERRA_UUSD_TOKEN_DECIMALS=6
BITSONG_TOKEN_DECIMALS=6
SECRET_TOKEN_DECIMALS=6
COMDEX_TOKEN_DECIMALS=6

for chain in AKASH SENTINEL IRIS PERSISTENCE CRYPTO_ORG REGEN OSMOSIS JUNO IXO LIKECOIN BITSONG BAND EMONEY EMONEY_EUR TERRA TERRA_UUSD SECRET COMDEX
do
  # do some munging
  chain_id=$chain"_CHAIN_ID"
  chain_id=${!chain_id}
  base_denom=$chain"_BASE_DENOM"
  base_denom=${!base_denom}
  channel_id=$chain"_CHANNEL_ID"
  channel_id=${!channel_id}
  counterparty_channel_id=$chain"_COUNTERPARTY_CHANNEL_ID"
  counterparty_channel_id=${!counterparty_channel_id}
  token_decimals=$chain"_TOKEN_DECIMALS"
  token_decimals=${!token_decimals}
  token_address=$chain"_TOKEN_ADDRESS"
  token_address=${!token_address}

  filename="$(tr [A-Z] [a-z] <<< "$chain")"

  # generate the IBC json
  sifnoded q tokenregistry generate -o json \
  	--token_base_denom=$base_denom \
  	--token_ibc_counterparty_chain_id=$chain_id \
        --token_ibc_channel_id=$channel_id \
        --token_ibc_counterparty_channel_id=$counterparty_channel_id \
  	--token_ibc_counterparty_denom="" \
  	--token_unit_denom="" \
  	--token_decimals=$token_decimals \
  	--token_display_name="" \
        --token_network=$chain_id \
        --token_address=$token_address \
  	--token_external_symbol="" \
  	--token_permission_clp=true \
  	--token_permission_ibc_export=true \
  	--token_permission_ibc_import=true | jq > $SIFCHAIN_ID/${filename}.json
  
  echo "\n\ngenerated entry for $chain $chain_id"
  cat $SIFCHAIN_ID/${filename}.json | jq
done

