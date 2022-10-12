#!/usr/bin/env bash

source ./data/margin-beta-users.sh

for addr in $users
do
  sifnoded tx margin whitelist $addr \
    --from=$ADMIN_KEY \
  	--gas=500000 \
  	--gas-prices=0.5rowan \
  	--chain-id $SIFCHAIN_ID \
  	--node $SIFNODE \
  	--broadcast-mode block \
  	--yes
done