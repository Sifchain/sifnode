#!/bin/bash 

addr=$1
shift
token=$1
shift

cd $SMART_CONTRACTS_DIR
yarn peggy:getTokenBalance ${addr:=${BRIDGE_BANK_ADDRESS}} ${token:=eth}
