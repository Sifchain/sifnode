#!/usr/bin/env bash

sifnoded tx tokenregistry register ./data/atom_all_permissions.json \
	--from $ADMIN_KEY \
	--gas=500000 \
	--gas-prices=0.5rowan \
	--chain-id $SIFCHAIN_ID \
	--node $SIFNODE \
	--broadcast-mode block \
	--yes