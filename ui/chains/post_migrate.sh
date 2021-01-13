#!/bin/bash

BASE_DIR=$PWD/../../

atk_address=$(cat $PWD/ethereum/build/contracts/AliceToken.json | jq -r '.networks["5777"].address') 
btk_address=$(cat $PWD/ethereum/build/contracts/BobToken.json | jq -r '.networks["5777"].address') 

if [[ -z "$atk_address" ]]; then 
  echo "Could not get atk address from json"
  exit 1
fi

if [[ -z "$btk_address" ]]; then 
  echo "Could not get btk address from json"
  exit 1
fi

# Whitelist test tokens
cd $BASE_DIR/smart-contracts

yarn peggy:whiteList "$atk_address" true
yarn peggy:whiteList "$btk_address" true

# Update local test addresses
cd $BASE_DIR/ui/core/
./scripts/updateLocalTestAddresses.js

