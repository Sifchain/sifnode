#!/bin/bash -x

sifnoded \
    q \
    tokenregistry \
    entries \
    --node ${SIF_NODE} \
    --chain-id=${SIF_CHAINID} \
    --output json \
    | jq '.entries[] | if .is_whitelisted == true then {denom: (if .base_denom != "" then .base_denom else .denom end), decimals: .decimals} else empty end'