#!/bin/bash

# This script should be run with a CWD that is the local folder
. ../credentials.sh
. ../../../smart-contracts/.env


BRIDGE_TOKEN_ADDRESS="0x82D50AD3C1091866E258Fd0f1a7cC9674609D254"

echo "ETHEREUM_PRIVATE_KEY=$ETHEREUM_PRIVATE_KEY"
echo "BRIDGE_TOKEN_ADDRESS=$BRIDGE_TOKEN_ADDRESS"
echo "SHADOWFIEND_NAME=$SHADOWFIEND_NAME"

ETHEREUM_PRIVATE_KEY=$ETHEREUM_PRIVATE_KEY ebrelayer init \
  tcp://localhost:26657 \
  ws://localhost:7545/ \
  "$BRIDGE_TOKEN_ADDRESS" \
  "$SHADOWFIEND_NAME" \
  "$SHADOWFIEND_MNEMONIC" \
  --chain-id=sifchain \
  --gas 300000 \
  --gas-adjustment 1.5