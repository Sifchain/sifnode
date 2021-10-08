#!/bin/bash -x

# sifnoded \
#     q \
#     ibc \
#     channel \
#     channels \
#     --chain-id=${SIF_CHAINID} \
#     --node ${SIF_NODE} \
#     --output json \
#     | jq '.channels[] | if .channel_id == "channel-61" then {src: .channel_id, dst: .counterparty.channel_id} else empty end'

sifnoded \
    q \
    tokenregistry \
    entries \
    --node ${SIF_NODE} \
    --chain-id=${SIF_CHAINID} \
    --output json \
    | jq '.entries[] | if .ibc_channel_id != "" then {denom: .base_denom, ibc_channel_id: .ibc_channel_id, ibc_counterparty_channel_id: .ibc_counterparty_channel_id, ibc_counterparty_chain_id: .ibc_counterparty_chain_id, permissions: .permissions} else empty end'