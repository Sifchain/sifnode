#!/bin/bash

. ./credentials.sh

if [[ -f "../../smart-contracts/.env" ]]; then
  . ../../smart-contracts/.env
fi

BASE_DIR=$PWD/../../

ATK_ADDRESS=$(cat $PWD/ethereum/build/contracts/AliceToken.json | jq -r '.networks["5777"].address') 
BTK_ADDRESS=$(cat $PWD/ethereum/build/contracts/BobToken.json | jq -r '.networks["5777"].address') 
USDC_ADDRESS=$(cat $PWD/ethereum/build/contracts/UsdCoin.json | jq -r '.networks["5777"].address') 
LINK_ADDRESS=$(cat $PWD/ethereum/build/contracts/LinkCoin.json | jq -r '.networks["5777"].address') 
BRIDGE_TOKEN_ADDRESS=$(cat $BASE_DIR/smart-contracts/build/contracts/BridgeToken.json | jq -r '.networks["5777"].address') 

if [[ -z "$ATK_ADDRESS" ]]; then 
  echo "Could not get atk address from json"
  exit 1
fi

if [[ -z "$BTK_ADDRESS" ]]; then 
  echo "Could not get btk address from json"
  exit 1
fi


if [[ -z "$BRIDGE_TOKEN_ADDRESS" ]]; then 
  echo "Could not get bridge token address from json"
  exit 1
fi

cd $BASE_DIR/smart-contracts

# Set token limits
UPDATE_ADDRESS=0x0000000000000000000000000000000000000000 npx truffle exec scripts/setTokenLockBurnLimit.js 31000000000000000000
UPDATE_ADDRESS=$BRIDGE_TOKEN_ADDRESS npx truffle exec scripts/setTokenLockBurnLimit.js 10000000000000000000000000

# Whitelist test tokens
yarn peggy:whiteList "$ATK_ADDRESS" true
yarn peggy:whiteList "$BTK_ADDRESS" true
yarn peggy:whiteList "$USDC_ADDRESS" true
yarn peggy:whiteList "$LINK_ADDRESS" true

# Update local test addresses
cd $BASE_DIR/ui/core/
./scripts/updateLocalTestAddresses.js

