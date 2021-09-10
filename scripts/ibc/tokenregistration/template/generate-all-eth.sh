#!/bin/sh

# ROWAN

sifnoded q tokenregistry generate \
	--token_base_denom=rowan \
	--token_decimals=18 \
	--token_unit_denom=rowan \
	--token_ibc_counterparty_denom=xrowan \
	--token_display_name="Rowan" \
	--token_external_symbol="eRowan" \
	--token_permission_clp=true \
	--token_permission_ibc_export=true \
	--token_permission_ibc_import=true > ./$SIFCHAIN_ID/rowan.json

sifnoded q tokenregistry generate \
	--token_base_denom=xrowan \
	--token_decimals=10 \
	--token_unit_denom=rowan \
	--token_display_name="Rowan" \
	--token_external_symbol="eRowan" \
	--token_permission_clp=false \
	--token_permission_ibc_export=false \
	--token_permission_ibc_import=true > ./$SIFCHAIN_ID/xrowan.json

# CETH

sifnoded q tokenregistry generate \
	--token_base_denom=ceth \
	--token_decimals=18 \
	--token_ibc_counterparty_denom=xeth \
	--token_display_name="ETH" \
	--token_external_symbol="ETH" \
	--token_permission_clp=true \
	--token_permission_ibc_export=true \
	--token_permission_ibc_import=true > ./$SIFCHAIN_ID/ceth.json

sifnoded q tokenregistry generate \
	--token_base_denom=xeth \
	--token_decimals=10 \
	--token_unit_denom=ceth \
	--token_display_name="ETH" \
	--token_external_symbol="ETH" \
	--token_permission_clp=false \
	--token_permission_ibc_export=false \
	--token_permission_ibc_import=true > ./$SIFCHAIN_ID/xeth.json

# CUSDC

sifnoded q tokenregistry generate \
	--token_base_denom=cusdc \
	--token_decimals=6 \
	--token_display_name="USDC" \
	--token_external_symbol="USDC" \
	--token_permission_clp=true \
	--token_permission_ibc_export=true \
	--token_permission_ibc_import=true > ./$SIFCHAIN_ID/cusdc.json

# CUSDT

sifnoded q tokenregistry generate \
	--token_base_denom=cusdt \
	--token_decimals=6 \
	--token_display_name="USDC" \
	--token_external_symbol="USDC" \
	--token_permission_clp=true \
	--token_permission_ibc_export=true \
	--token_permission_ibc_import=true > ./$SIFCHAIN_ID/cusdt.json
