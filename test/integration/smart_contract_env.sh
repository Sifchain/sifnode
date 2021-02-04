#!/bin/bash

# Usage:
#
# $0 directory_name
#
# where directory_name points to a directory with the json for our smart contracts

. $(dirname $0)/shell_utilities.sh

json_contracts_dir=${1}
shift

echo export BRIDGE_REGISTRY_ADDRESS=$(cat $json_contracts_dir/BridgeRegistry.json | jq -r '.networks["3"].address')
echo export BRIDGE_TOKEN_ADDRESS=$(cat $json_contracts_dir/BridgeToken.json | jq -r '.networks["3"].address')
echo export BRIDGE_BANK_ADDRESS=$(cat $json_contracts_dir/BridgeBank.json | jq -r '.networks["3"].address')