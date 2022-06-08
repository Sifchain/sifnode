#!/bin/sh

# sh ./generate-erc20-jsons.sh testnet

. ./envs/$1.sh

echo "\n\ngenerating and storing all entries for network $SIFCHAIN_ID"

mkdir -p ./$SIFCHAIN_ID

sifnoded q tokenregistry generate \
	--token_base_denom=cxft \
	--token_denom=cxft \
	--token_decimals=18 \
	--token_display_name="Offshift" \
	--token_external_symbol="XFT" \
	--token_permission_clp=true \
	--token_permission_ibc_export=true \
	--token_permission_ibc_import=true \
	-o json | jq > $SIFCHAIN_ID/cxft.json

echo "\n\ngenerated entry for cxft"

sifnoded q tokenregistry generate \
	--token_base_denom=cuos \
	--token_denom=cuos \
	--token_decimals=4 \
	--token_display_name="Ultra Token" \
	--token_external_symbol="UOS" \
	--token_permission_clp=true \
	--token_permission_ibc_export=true \
	--token_permission_ibc_import=true -o json | jq > $SIFCHAIN_ID/cuos.json

echo "\n\ngenerated entry for cuos"

sifnoded q tokenregistry generate \
	--token_denom=cnewo \
	--token_base_denom=cnewo \
	--token_decimals=18 \
	--token_display_name="New Order" \
	--token_external_symbol="NEWO" \
	--token_permission_clp=true \
	--token_permission_ibc_export=true \
	--token_permission_ibc_import=true -o json | jq > $SIFCHAIN_ID/cnewo.json

echo "\n\ngenerated entry for cnewo"

sifnoded q tokenregistry generate \
	--token_denom=cosqth \
	--token_base_denom=cosqth \
	--token_decimals=18 \
	--token_display_name="Opyn Squeeth" \
	--token_external_symbol="oSQTH" \
	--token_permission_clp=true \
	--token_permission_ibc_export=true \
	--token_permission_ibc_import=true -o json | jq > $SIFCHAIN_ID/cosqth.json

echo "\n\ngenerated entry for cosqth"

sifnoded q tokenregistry generate \
	--token_denom=cgala \
	--token_base_denom=cgala \
	--token_decimals=8 \
	--token_display_name="Gala" \
	--token_external_symbol="GALA" \
	--token_permission_clp=true \
	--token_permission_ibc_export=true \
	--token_permission_ibc_import=true -o json | jq > $SIFCHAIN_ID/cgala.json

echo "\n\ngenerated entry for cgala"


sifnoded q tokenregistry generate \
	--token_denom=cpush \
	--token_base_denom=cpush \
	--token_decimals=18 \
	--token_display_name="Ethereum Push Notification Service" \
	--token_external_symbol="PUSH" \
	--token_permission_clp=true \
	--token_permission_ibc_export=true \
	--token_permission_ibc_import=true -o json | jq > $SIFCHAIN_ID/cpush.json

echo "\n\ngenerated entry for cpush"


sifnoded q tokenregistry generate \
	--token_denom=cmc \
	--token_base_denom=cmc \
	--token_decimals=18 \
	--token_display_name="Merit Circle" \
	--token_external_symbol="MC" \
	--token_permission_clp=true \
	--token_permission_ibc_export=true \
	--token_permission_ibc_import=true -o json | jq > $SIFCHAIN_ID/cmc.json

echo "\n\ngenerated entry for cmc"

sifnoded q tokenregistry generate \
	--token_denom=cinj \
	--token_base_denom=cinj \
	--token_decimals=18 \
	--token_display_name="Injective Token" \
	--token_external_symbol="INJ" \
	--token_permission_clp=true \
	--token_permission_ibc_export=true \
	--token_permission_ibc_import=true -o json | jq > $SIFCHAIN_ID/cinj.json

echo "\n\ngenerated entry for cinj"

sifnoded q tokenregistry generate \
	--token_denom=ccudos \
	--token_base_denom=ccudos \
	--token_decimals=18 \
	--token_display_name="Cudos" \
	--token_external_symbol="CUDOS" \
	--token_permission_clp=true \
	--token_permission_ibc_export=true \
	--token_permission_ibc_import=true -o json | jq > $SIFCHAIN_ID/ccudos.json

echo "\n\ngenerated entry for ccudos"
