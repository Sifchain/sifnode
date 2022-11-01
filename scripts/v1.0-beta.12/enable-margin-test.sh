#!/usr/bin/env bash

sifnoded tx margin update-pools ./data/temp_pools.json \
	--closed-pools ./data/closed_pools.json \
  --from=$ADMIN_KEY \
	--gas=500000 \
	--gas-prices=0.5rowan \
	--chain-id $SIFCHAIN_ID \
	--node $SIFNODE \
	--broadcast-mode block \
	--yes

sifnoded tx margin whitelist sif1mwmrarhynjuau437d07p42803rntfxqjun3pfu \
  --from=$ADMIN_KEY \
	--gas=500000 \
	--gas-prices=0.5rowan \
	--chain-id $SIFCHAIN_ID \
	--node $SIFNODE \
	--broadcast-mode block \
	--yes