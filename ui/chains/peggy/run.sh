#!/bin/bash

# This script should be run with a CWD that is the local folder
. ../credentials.sh
. ../../../smart-contracts/.env


BRIDGE_TOKEN_ADDRESS="0x82D50AD3C1091866E258Fd0f1a7cC9674609D254"

echo "  "
echo "-----------------------------------------------------"
echo "BRIDGE_TOKEN_ADDRESS='$BRIDGE_TOKEN_ADDRESS'"
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

ETHEREUM_PRIVATE_KEY=$ETHEREUM_PRIVATE_KEY ebrelayer init \
  tcp://localhost:26657 \
  ws://localhost:7545/ \
  "$BRIDGE_TOKEN_ADDRESS" \
  "$SHADOWFIEND_NAME" \
  "$SHADOWFIEND_MNEMONIC" \
  --chain-id=sifchain \
  --gas 300000 \
  --gas-adjustment 1.5