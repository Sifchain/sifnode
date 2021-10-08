#!/bin/bash -x

export SIF_CHAINID="sifchain-testnet-042-ibc"
export SIF_NODE="http://rpc-testnet-042-ibc.sifchain.finance:80"
export COSMOS_CHAINID="cosmoshub-testnet"
export COSMOS_NODE="https://rpc.testnet.cosmos.network:443"
export SIFTOCOSMOS_CHANNEL_ID="channel-0"
export COSMOSTOSIF_CHANNEL_ID="channel-86"
export SIF_CHANNEL_IDS=channel-0,channel-1,channel-2,channel-3,channel-4
export COSMOS_CHANNEL_IDS=channel-86,channel-12,channel-16,channel-5

export CHAINS_BINARY="gaiad"
export CHAINS_NODE="https://rpc.testnet.cosmos.network:443"
export CHAINS_ID="cosmoshub-testnet"
export CHAINS_DENOM="uphoton"
export CHAINS_FEES="0"
export SIFTOCHAINS_CHANNEL_ID="channel-23"
export CHAINSTOSIF_CHANNEL_ID="channel-16"