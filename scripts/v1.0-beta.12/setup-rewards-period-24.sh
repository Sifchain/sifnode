#!/usr/bin/env bash

sifnoded tx clp reward-period --path=./data/rp_24.json \
	--from $ADMIN_KEY \
	--gas=500000 \
	--gas-prices=0.5rowan \
	--chain-id $SIFCHAIN_ID \
	--node $SIFNODE \
	--broadcast-mode block \
	--yes