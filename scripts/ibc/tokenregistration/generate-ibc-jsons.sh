#!/bin/sh

# sh ./generate-ibc-jsons.sh testnet

if [ $# -ne 1 ]; then
  echo "Need to pass in the environment. EX: generate-ibc-jsons.sh devnet"
  exit 1
fi

. ./envs/$1.sh 

echo "\n\ngenerating and storing all entries for network $SIFCHAIN_ID"

mkdir -p ./$SIFCHAIN_ID

for chain in AKASH SENTINEL IRIS PERSISTENCE CRYPTO_ORG REGEN OSMOSIS JUNO IXO LIKECOIN BITSONG BAND EMONEY EMONEY_EEUR TERRA TERRA_UUSD SECRET COMDEX
do
  # do some munging
  chain_id=$chain"_CHAIN_ID"
  chain_id=${!chain_id}

  # skip the empties
  if [ -z "$chain_id" ]; then
    continue
  fi

  base_denom=$chain"_BASE_DENOM"
  base_denom=${!base_denom}
  channel_id=$chain"_CHANNEL_ID"
  channel_id=${!channel_id}
  counterparty_channel_id=$chain"_COUNTERPARTY_CHANNEL_ID"
  counterparty_channel_id=${!counterparty_channel_id}

  # almost always they are 6 but this is the place to override them
  token_decimals=0
  case $chain in
    LIKECOIN)
      token_decimals=9
      ;;
    SECRET)
      token_decimals=18
      ;;
    CRYPTO_ORG)
      token_decimals=8
      ;;
    COMDEX)
      token_decimals=18
      ;;
    *)
      token_decimals=6
      ;;
  esac

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

