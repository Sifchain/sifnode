#!/bin/bash -x

export SIF_CHAINID="sifchain-devnet-1"
export SIF_NODE="https://rpc-devnet.sifchain.finance:443"
export COSMOS_CHAINID="cosmoshub-testnet"
export COSMOS_NODE="https://rpc.testnet.cosmos.network:443"
export SIFTOCOSMOS_CHANNEL_ID="channel-61"
export COSMOSTOSIF_CHANNEL_ID="channel-128"

# export CHAINS_BINARY="gaiad,iris,akash"
# export CHAINS_NODE="https://rpc.testnet.cosmos.network:443,http://35.234.10.84:26657,http://147.75.32.35:26657"
# export CHAINS_ID="cosmoshub-testnet,nyancat-8,akash-testnet-6"
# export CHAINS_DENOM="uphoton,unyan,uakt"
# export CHAINS_FEES="5000"
# export SIFTOCHAINS_CHANNEL_ID="channel-61,channel-42,channel-77"
# export CHAINSTOSIF_CHANNEL_ID="channel-128,channel-13,channel-44"

export CHAINS_BINARY="iris"
export CHAINS_NODE="http://35.234.10.84:26657"
export CHAINS_ID="nyancat-8"
export CHAINS_DENOM="unyan"
export CHAINS_FEES="400"
export SIFTOCHAINS_CHANNEL_ID="channel-117"
export CHAINSTOSIF_CHANNEL_ID="channel-24"

# export CHAINS_BINARY="osmosis"
# export CHAINS_NODE="http://osmosis-rpc.osmosis-node:26657"
# export CHAINS_ID="osmosis-1"
# export CHAINS_DENOM="uosmo"
# export CHAINS_FEES="0"
# export SIFTOCHAINS_CHANNEL_ID="channel-74"
# export CHAINSTOSIF_CHANNEL_ID="channel-22"

# export CHAINS_BINARY="akash"
# export CHAINS_NODE="http://147.75.32.35:26657"
# export CHAINS_ID="akash-testnet-6"
# export CHAINS_DENOM="uakt"
# export CHAINS_FEES="5000"
# export SIFTOCHAINS_CHANNEL_ID="channel-115"
# export CHAINSTOSIF_CHANNEL_ID="channel-65"

# export CHAINS_BINARY="sentinelhub"
# export CHAINS_NODE="http://rpc.sentinel.co:26657"
# export CHAINS_ID="sentinelhub-2"
# export CHAINS_DENOM="udvpn"
# export CHAINS_FEES="20000"
# export SIFTOCHAINS_CHANNEL_ID="channel-79"
# export CHAINSTOSIF_CHANNEL_ID="channel-19"

# export CHAINS_BINARY="gaiad"
# export CHAINS_NODE="https://rpc.testnet.cosmos.network:443"
# export CHAINS_ID="cosmoshub-testnet"
# export CHAINS_DENOM="uphoton"
# export CHAINS_FEES="0"
# export SIFTOCHAINS_CHANNEL_ID="channel-101"
# export CHAINSTOSIF_CHANNEL_ID="channel-3"

# export CHAINS_BINARY="persistenceCore"
# export CHAINS_NODE="https://persistence.testnet.rpc.audit.one:443"
# export CHAINS_ID="test-core-1"
# export CHAINS_DENOM="uxprt"
# export CHAINS_FEES="0"
# export SIFTOCHAINS_CHANNEL_ID="channel-120"
# export CHAINSTOSIF_CHANNEL_ID="channel-25"

export INFURA_PROJECT_ID=c413023ff7944d21b694664b31a52faf
export OPERATOR_ADDRESS=0x1e0220B251eE648C7F3B6Fc31E6d309141f2e464
export OPERATOR_PRIVATE_KEY=30dd94b42b731aa5fe738353d897fb938cf7e2a1dbce629dead1b3294ede4f3c

export ETH_CHAINID="3"
export ETH_ADDRESS="0x5171050beb52148aB834Fb21E3E30FA429470c46"
export ETH_PRIVATE_KEY="4ef63770f02888abded959c550e6d1060859bad0f344abc78009af37db936d6d"
export ETH_NETWORK="ropsten"
export BRIDGEBANK_ADDRESS="0xB75849afEF2864977a858073458Cb13F9410f8e5"
export DEPLOYMENT_NAME="sifchain-devnet-1"

export ETHEREUM_NETWORK=ropsten
export ETHEREUM_NETWORK_ID=3

# export BASEDIR=/sifnode

# export SMART_CONTRACTS_DIR=$BASEDIR/smart-contracts
# export SOLIDITY_JSON_PATH=$BASEDIR/smart-contracts/deployments/$DEPLOYMENT_NAME
# export SMART_CONTRACT_ARTIFACT_DIR=$SOLIDITY_JSON_PATH

# export BRIDGE_REGISTRY_ADDRESS=$(cat $SOLIDITY_JSON_PATH/BridgeRegistry.json | jq -r ".networks[\"$ETHEREUM_NETWORK_ID\"].address")
# export BRIDGE_TOKEN_ADDRESS=$(cat $SOLIDITY_JSON_PATH/BridgeToken.json | jq -r ".networks[\"$ETHEREUM_NETWORK_ID\"].address")
# export BRIDGE_BANK_ADDRESS=$(cat $SOLIDITY_JSON_PATH/BridgeBank.json | jq -r ".networks[\"$ETHEREUM_NETWORK_ID\"].address")