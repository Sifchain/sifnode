#!/bin/bash

# This script should be run with a CWD that is the local folder
. $PWD/../credentials.sh

if [[ -f "$PWD/../../../smart-contracts/.env" ]]; then
  . $PWD/../../../smart-contracts/.env
fi

# Required to run ebrelayer
export BRIDGE_TOKEN_ADDRESS=$(cat $PWD/../../../smart-contracts/build/contracts/BridgeToken.json | jq -r '.networks["5777"].address')
export BRIDGE_REGISTRY_ADDRESS=$(cat $PWD/../../../smart-contracts/build/contracts/BridgeRegistry.json | jq -r '.networks["5777"].address') 

echo "  "
echo "-----------------------------------------------------"
echo "BRIDGE_TOKEN_ADDRESS='$BRIDGE_TOKEN_ADDRESS'"
echo "BRIDGE_REGISTRY_ADDRESS='$BRIDGE_REGISTRY_ADDRESS'"
echo "CI=$CI"
echo "CONSENSUS_THRESHOLD='$CONSENSUS_THRESHOLD'"
echo "EROWAN_ADDRESS='$EROWAN_ADDRESS'"
echo "ETHEREUM_PRIVATE_KEY='$ETHEREUM_PRIVATE_KEY'"
echo "INFURA_PROJECT_ID='$INFURA_PROJECT_ID'"
echo "INITIAL_VALIDATOR_ADDRESSES='$INITIAL_VALIDATOR_ADDRESSES'"
echo "INITIAL_VALIDATOR_POWERS='$INITIAL_VALIDATOR_POWERS'"
echo "MAINNET_GAS_PRICE='$MAINNET_GAS_PRICE'"
echo "MNEMONIC='$MNEMONIC'"
echo "OPERATOR='$OPERATOR'"
echo "OWNER='$OWNER'"
echo "SHADOWFIEND_NAME='$SHADOWFIEND_NAME'"
echo "-----------------------------------------------------"
echo "  "


cd $BASE_DIR && ETHEREUM_PRIVATE_KEY=$ETHEREUM_PRIVATE_KEY ebrelayer init \
  tcp://localhost:26657 \
  ws://localhost:7545/ \
  "$BRIDGE_REGISTRY_ADDRESS" \
  "$SHADOWFIEND_NAME" \
  "$SHADOWFIEND_MNEMONIC" \
  --chain-id=sifchain-local \
  --gas 300000 \
  --gas-adjustment 1.5